let nickname = document.getElementById("nickname");
nickname.value = localStorage.nickname || "";
nickname.oninput = (e) => {
    localStorage.nickname = e.target.value;
};

let postOrder = document.getElementById("postOrder");
postOrder.value = localStorage.postOrder || "chronological";
postOrder.onchange = (e) => {
    localStorage.postOrder = e.target.value;
};

let norlogeStyle = document.getElementById("norlogeStyle");
norlogeStyle.value = localStorage.norlogeStyle || "auto";
norlogeStyle.onchange = (e) => {
    localStorage.norlogeStyle = e.target.value;
};

let backButton = document.getElementById("backButton");
backButton.onclick = () => {
    window.location.href = "/";
}