class Gc2Alarm extends HTMLElement{

    constructor() {
        super();
    }

    connectedCallback() {

        setInterval(() => {
            let dateLocal = new Date();
            let hours = dateLocal.getHours();
            let minutes = dateLocal.getMinutes();
            if(minutes !== 59) {
                this.toggleAttribute("preums", false);
                this.innerText = `${hours < 10 ? '0':''}${hours}:${minutes < 10 ? '0':''}${minutes}`;
            } else {
                this.toggleAttribute("preums", true);
                let seconds = dateLocal.getSeconds();
                this.innerText = `${hours < 10 ? '0':''}${hours}:${minutes < 10 ? '0':''}${minutes}:${seconds < 10 ? '0':''}${seconds}`;
            }
        }, 1000);

    }
}

customElements.define('gc2-alarm', Gc2Alarm);