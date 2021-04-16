let nickname = document.getElementById("nickname");
nickname.value = localStorage.nickname || "";
nickname.oninput = (e) => {
    localStorage.nickname = e.target.value;
};

let backButton = document.getElementById("backButton");
backButton.onclick = () => {
    window.location.href = "/";
}