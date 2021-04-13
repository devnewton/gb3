/**
 * @callback Gc2TribuneNavigator~navigateCallback
 * @param {string} tribuneName
 */

 class Gc2TribuneNavigator extends HTMLElement {
    constructor() {
        super();
    }

    /**
     * @type Gc2TribuneNavigator~navigateCallback
     */
    onnavigate;

    addTribune(tribuneName) {
        let tribuneButton = document.createElement("li");
        tribuneButton.innerText = tribuneName;
        tribuneButton.onclick = () => {
            if(this.onnavigate) {
                this.querySelectorAll("li").forEach((element) => {
                    element.classList.toggle("gc2-tribune-navigator-active", false);
                });
                tribuneButton.classList.toggle("gc2-tribune-navigator-active", true);
                this.onnavigate(tribuneName);
            }
        }
        this.appendChild(tribuneButton);
    }
}

customElements.define('gc2-tribune-navigator', Gc2TribuneNavigator);