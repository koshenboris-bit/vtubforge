import * as THREE from "three";
import { GLTFLoader } from "three/addons/loaders/GLTFLoader.js";
import { OrbitControls } from "three/addons/controls/OrbitControls.js";
import { VRMLoaderPlugin, VRMUtils } from "@pixiv/three-vrm";
import { FaceLandmarker, FilesetResolver } from "https://cdn.jsdelivr.net/npm/@mediapipe/tasks-vision@0.10.12";

const canvas = document.getElementById("avatarCanvas");
const modelInput = document.getElementById("modelInput");
const cameraButton = document.getElementById("cameraButton");
const resetButton = document.getElementById("resetButton");
const video = document.getElementById("webcamVideo");
const stageEmpty = document.getElementById("stageEmpty");
const trackingStatus = document.getElementById("trackingStatus");
const faceStatus = document.getElementById("faceStatus");
const modelStatus = document.getElementById("modelStatus");

const clock = new THREE.Clock();
const scene = new THREE.Scene();
scene.fog = new THREE.Fog(0x070b18, 7, 18);

const camera = new THREE.PerspectiveCamera(32, 1, 0.1, 100);
camera.position.set(0, 1.35, 5.2);

const renderer = new THREE.WebGLRenderer({ canvas, antialias: true, alpha: true });
renderer.setPixelRatio(Math.min(window.devicePixelRatio || 1, 2));
renderer.outputColorSpace = THREE.SRGBColorSpace;

const controls = new OrbitControls(camera, canvas);
controls.target.set(0, 1.15, 0);
controls.enableDamping = true;
controls.enablePan = false;
controls.minDistance = 2.4;
controls.maxDistance = 8;

const keyLight = new THREE.DirectionalLight(0xffffff, 2.6);
keyLight.position.set(2.5, 4, 4);
scene.add(keyLight);
scene.add(new THREE.HemisphereLight(0x8bdcff, 0x27124d, 1.8));

const floor = new THREE.Mesh(
  new THREE.CircleGeometry(2.2, 96),
  new THREE.MeshBasicMaterial({ color: 0x38bdf8, transparent: true, opacity: 0.08 })
);
floor.rotation.x = -Math.PI / 2;
floor.position.y = -1.08;
scene.add(floor);

const loader = new GLTFLoader();
loader.register((parser) => new VRMLoaderPlugin(parser));

let avatarScene = null;
let currentVrm = null;
let headBone = null;
let neckBone = null;
let morphMeshes = [];
let faceLandmarker = null;
let cameraRunning = false;
let lastVideoTime = -1;
let targetPose = { x: 0, y: 0, z: 0 };
let targetExpressions = { blink: 0, smile: 0, mouth: 0, lookLeft: 0, lookRight: 0 };
let loadedObjectUrl = null;

function tx(key) {
  return window.t ? window.t(key) : key;
}

function setText(element, key) {
  if (element) element.textContent = tx(key);
}

function resizeRenderer() {
  const rect = canvas.parentElement.getBoundingClientRect();
  renderer.setSize(rect.width, rect.height, false);
  camera.aspect = rect.width / Math.max(rect.height, 1);
  camera.updateProjectionMatrix();
}

function disposeNode(root) {
  root.traverse((node) => {
    if (node.geometry) node.geometry.dispose?.();
    const materials = Array.isArray(node.material) ? node.material : node.material ? [node.material] : [];
    materials.forEach((material) => material.dispose?.());
  });
}

function clearAvatar() {
  if (avatarScene) {
    scene.remove(avatarScene);
    disposeNode(avatarScene);
  }
  currentVrm = null;
  avatarScene = null;
  headBone = null;
  neckBone = null;
  morphMeshes = [];
}

function collectMorphTargets(root) {
  morphMeshes = [];
  root.traverse((node) => {
    if (node.isMesh && node.morphTargetDictionary && node.morphTargetInfluences) {
      morphMeshes.push(node);
    }
  });
}

function fitModel(root) {
  const box = new THREE.Box3().setFromObject(root);
  const size = new THREE.Vector3();
  const center = new THREE.Vector3();
  box.getSize(size);
  box.getCenter(center);

  root.position.sub(center);
  const maxAxis = Math.max(size.x, size.y, size.z) || 1;
  root.scale.setScalar(2.65 / maxAxis);
  root.position.y -= 0.1;
}

function setupVrm(vrm) {
  currentVrm = vrm;
  avatarScene = vrm.scene;

  VRMUtils.removeUnnecessaryVertices(avatarScene);
  VRMUtils.removeUnnecessaryJoints(avatarScene);
  avatarScene.traverse((object) => {
    object.frustumCulled = false;
  });

  fitModel(avatarScene);
  scene.add(avatarScene);

  headBone = vrm.humanoid?.getNormalizedBoneNode("head") || null;
  neckBone = vrm.humanoid?.getNormalizedBoneNode("neck") || null;
  collectMorphTargets(avatarScene);
}

function setupGenericGltf(gltf) {
  avatarScene = gltf.scene;
  fitModel(avatarScene);
  collectMorphTargets(avatarScene);
  scene.add(avatarScene);
}

function setMorphByPattern(patterns, value) {
  for (const mesh of morphMeshes) {
    const dict = mesh.morphTargetDictionary;
    for (const [name, index] of Object.entries(dict)) {
      if (patterns.some((pattern) => pattern.test(name))) {
        mesh.morphTargetInfluences[index] = THREE.MathUtils.lerp(mesh.morphTargetInfluences[index] || 0, value, 0.45);
      }
    }
  }
}

function setVrmExpression(names, value) {
  const manager = currentVrm?.expressionManager;
  if (!manager) return false;

  let applied = false;
  for (const name of names) {
    try {
      manager.setValue(name, value);
      applied = true;
    } catch (_) {}
  }
  return applied;
}

function applyExpressions() {
  const blink = THREE.MathUtils.clamp(targetExpressions.blink, 0, 1);
  const smile = THREE.MathUtils.clamp(targetExpressions.smile, 0, 1);
  const mouth = THREE.MathUtils.clamp(targetExpressions.mouth, 0, 1);

  if (currentVrm?.expressionManager) {
    setVrmExpression(["blink", "blinkLeft", "blinkRight"], blink);
    setVrmExpression(["happy", "relaxed"], smile * 0.85);
    setVrmExpression(["aa", "oh"], mouth);
    currentVrm.expressionManager.update?.();
    return;
  }

  setMorphByPattern([/blink/i, /eye.*close/i, /close.*eye/i], blink);
  setMorphByPattern([/smile/i, /happy/i, /joy/i, /fun/i], smile);
  setMorphByPattern([/jaw.*open/i, /mouth.*open/i, /^aa$/i, /^a$/i, /v_aa/i], mouth);
}

function blendshapeScore(faceBlendshapes, names) {
  const categories = faceBlendshapes?.[0]?.categories || [];
  let best = 0;
  for (const category of categories) {
    if (names.includes(category.categoryName)) best = Math.max(best, category.score || 0);
  }
  return best;
}

function setPoseFromTransformation(results) {
  const matrixData = results.facialTransformationMatrixes?.[0]?.data;
  if (!matrixData) return false;

  const matrix = new THREE.Matrix4().fromArray(matrixData);
  const euler = new THREE.Euler().setFromRotationMatrix(matrix, "YXZ");
  targetPose.x = THREE.MathUtils.clamp(-euler.x * 0.55, -0.34, 0.34);
  targetPose.y = THREE.MathUtils.clamp(euler.y * 0.75, -0.55, 0.55);
  targetPose.z = THREE.MathUtils.clamp(-euler.z * 0.55, -0.32, 0.32);
  return true;
}

function setPoseFromLandmarks(landmarks) {
  const nose = landmarks[1] || landmarks[4] || landmarks[0];
  if (!nose) return;
  targetPose.y = THREE.MathUtils.clamp(-(nose.x - 0.5) * 1.45, -0.48, 0.48);
  targetPose.x = THREE.MathUtils.clamp(-(nose.y - 0.5) * 0.9, -0.32, 0.32);
  targetPose.z = 0;
}

function updateFromFace(results) {
  const landmarks = results.faceLandmarks?.[0];
  if (!landmarks) {
    setText(faceStatus, "tryFaceSearching");
    targetPose.x *= 0.88;
    targetPose.y *= 0.88;
    targetPose.z *= 0.88;
    targetExpressions.blink *= 0.8;
    targetExpressions.smile *= 0.82;
    targetExpressions.mouth *= 0.82;
    return;
  }

  setText(faceStatus, "tryFaceOn");
  if (!setPoseFromTransformation(results)) setPoseFromLandmarks(landmarks);

  const leftBlink = blendshapeScore(results.faceBlendshapes, ["eyeBlinkLeft"]);
  const rightBlink = blendshapeScore(results.faceBlendshapes, ["eyeBlinkRight"]);
  const leftSmile = blendshapeScore(results.faceBlendshapes, ["mouthSmileLeft"]);
  const rightSmile = blendshapeScore(results.faceBlendshapes, ["mouthSmileRight"]);
  const jawOpen = blendshapeScore(results.faceBlendshapes, ["jawOpen"]);
  const mouthFunnel = blendshapeScore(results.faceBlendshapes, ["mouthFunnel"]);

  targetExpressions.blink = Math.max(leftBlink, rightBlink);
  targetExpressions.smile = Math.max(leftSmile, rightSmile);
  targetExpressions.mouth = Math.max(jawOpen, mouthFunnel * 0.7);
}

async function initFaceTracker() {
  if (faceLandmarker) return;
  setText(trackingStatus, "tryStatusLoadingTracker");
  const fileset = await FilesetResolver.forVisionTasks("https://cdn.jsdelivr.net/npm/@mediapipe/tasks-vision@0.10.12/wasm");
  faceLandmarker = await FaceLandmarker.createFromOptions(fileset, {
    baseOptions: {
      modelAssetPath: "https://storage.googleapis.com/mediapipe-models/face_landmarker/face_landmarker/float16/latest/face_landmarker.task",
      delegate: "GPU"
    },
    runningMode: "VIDEO",
    numFaces: 1,
    outputFaceBlendshapes: true,
    outputFacialTransformationMatrixes: true
  });
}

async function startCamera() {
  try {
    await initFaceTracker();
    const stream = await navigator.mediaDevices.getUserMedia({ video: { width: 640, height: 480, facingMode: "user" }, audio: false });
    video.srcObject = stream;
    await video.play();
    cameraRunning = true;
    setText(trackingStatus, "tryStatusTracking");
    cameraButton.classList.add("hidden");
  } catch (error) {
    console.error(error);
    setText(trackingStatus, "tryStatusCameraError");
    setText(faceStatus, "tryFaceOff");
  }
}

function detectFace() {
  if (!cameraRunning || !faceLandmarker || video.readyState < 2) return;
  if (video.currentTime === lastVideoTime) return;
  lastVideoTime = video.currentTime;
  updateFromFace(faceLandmarker.detectForVideo(video, performance.now()));
}

function applyPose() {
  if (headBone) {
    headBone.rotation.x = THREE.MathUtils.lerp(headBone.rotation.x, targetPose.x, 0.18);
    headBone.rotation.y = THREE.MathUtils.lerp(headBone.rotation.y, targetPose.y, 0.18);
    headBone.rotation.z = THREE.MathUtils.lerp(headBone.rotation.z, targetPose.z, 0.18);
    if (neckBone) {
      neckBone.rotation.x = THREE.MathUtils.lerp(neckBone.rotation.x, targetPose.x * 0.28, 0.14);
      neckBone.rotation.y = THREE.MathUtils.lerp(neckBone.rotation.y, targetPose.y * 0.28, 0.14);
      neckBone.rotation.z = THREE.MathUtils.lerp(neckBone.rotation.z, targetPose.z * 0.18, 0.14);
    }
  } else if (avatarScene) {
    avatarScene.rotation.x = THREE.MathUtils.lerp(avatarScene.rotation.x, targetPose.x * 0.45, 0.12);
    avatarScene.rotation.y = THREE.MathUtils.lerp(avatarScene.rotation.y, targetPose.y * 0.45, 0.12);
    avatarScene.rotation.z = THREE.MathUtils.lerp(avatarScene.rotation.z, targetPose.z * 0.35, 0.12);
  }
}

function resetPose() {
  targetPose = { x: 0, y: 0, z: 0 };
  targetExpressions = { blink: 0, smile: 0, mouth: 0, lookLeft: 0, lookRight: 0 };
  if (headBone) headBone.rotation.set(0, 0, 0);
  if (neckBone) neckBone.rotation.set(0, 0, 0);
  if (avatarScene && !headBone) avatarScene.rotation.set(0, 0, 0);
}

async function loadModel(file) {
  clearAvatar();
  if (loadedObjectUrl) URL.revokeObjectURL(loadedObjectUrl);
  loadedObjectUrl = URL.createObjectURL(file);
  setText(trackingStatus, "tryStatusLoadingModel");
  modelStatus.textContent = file.name;

  try {
    const gltf = await loader.loadAsync(loadedObjectUrl);
    if (gltf.userData?.vrm) setupVrm(gltf.userData.vrm);
    else setupGenericGltf(gltf);

    stageEmpty.classList.add("hidden");
    const hasFaceExpressions = !!currentVrm?.expressionManager || morphMeshes.length > 0;
    setText(trackingStatus, hasFaceExpressions ? "tryStatusModelReadyFace" : "tryStatusModelReadyBasic");
  } catch (error) {
    console.error(error);
    setText(trackingStatus, "tryStatusModelError");
    setText(modelStatus, "tryModelEmpty");
    stageEmpty.classList.remove("hidden");
  }
}

function animate() {
  requestAnimationFrame(animate);
  const delta = clock.getDelta();
  detectFace();
  applyPose();
  applyExpressions();
  currentVrm?.update?.(delta);
  floor.rotation.z += 0.003;
  controls.update();
  renderer.render(scene, camera);
}

modelInput.addEventListener("change", (event) => {
  const file = event.target.files?.[0];
  if (file) loadModel(file);
});

cameraButton.addEventListener("click", startCamera);
resetButton.addEventListener("click", resetPose);
window.addEventListener("resize", resizeRenderer);
window.addEventListener("vtuberforge:languagechange", () => {
  if (!avatarScene) setText(modelStatus, "tryModelEmpty");
  if (!cameraRunning) setText(faceStatus, "tryFaceOff");
  if (!avatarScene) setText(trackingStatus, "tryStatusIdle");
});

resizeRenderer();
setText(trackingStatus, "tryStatusIdle");
setText(faceStatus, "tryFaceOff");
setText(modelStatus, "tryModelEmpty");
animate();
