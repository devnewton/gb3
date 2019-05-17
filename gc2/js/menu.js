class Gc2Menu extends HTMLElement {
    constructor() {
        super();
        this.style.display = "none";

        let totozButton = document.createElement('button');
        totozButton.innerText = "Totoz";
        this.appendChild(totozButton);

        let emojiButton = document.createElement('button');
        emojiButton.innerText = "Emoji";
        this.appendChild(emojiButton);

        let binButton = document.createElement('button');
        binButton.innerText = "Bin";
        this.appendChild(binButton);

        let settingsLink = document.createElement('a');
        settingsLink.innerText = "Settings";
        settingsLink.href = "/settings.html";
        this.appendChild(settingsLink);

        let closeButton = document.createElement('button');
        closeButton.innerText = "Close";
        closeButton.onclick = () => {
            this.style.display = "none";
            document.querySelector("gc2-main").style.display = "flex";
        }
        this.appendChild(closeButton);
    }
}
customElements.define('gc2-menu', Gc2Menu);