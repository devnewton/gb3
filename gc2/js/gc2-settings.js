let nickname = document.getElementById("nickname");
nickname.value = localStorage.nickname || "";
nickname.oninput = (e) => {
    localStorage.nickname = e.target.value;
};

let postOrder = document.getElementById("postOrder");
postOrder.value = localStorage.postOrder || "reverse-chronological";
postOrder.onchange = (e) => {
    localStorage.postOrder = e.target.value;
};

let norlogeStyle = document.getElementById("norlogeStyle");
norlogeStyle.value = localStorage.norlogeStyle || "longlong";
norlogeStyle.onchange = (e) => {
    localStorage.norlogeStyle = e.target.value;
};