class Gc2TribuneNavigatorButton extends HTMLElement{

    /**
     * @type string
     */
    tribuneName;

    constructor() {
        super();
    }

}

customElements.define('gc2-tribune-navigator-button', Gc2TribuneNavigatorButton);

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
        let tribuneButton = document.createElement("gc2-tribune-navigator-button");
        tribuneButton.tribuneName = tribuneName;
        tribuneButton.innerText = tribuneName;
        tribuneButton.onclick = () => {
            if(this.onnavigate) {
                this.onnavigate(tribuneName);
            }
        }
        this.appendChild(tribuneButton);
    }

    setActiveTribune(tribuneName) {
        this.querySelectorAll("gc2-tribune-navigator-button").forEach((element) => {
            element.toggleAttribute("active", element.tribuneName === tribuneName);
        });
    }
}

customElements.define('gc2-tribune-navigator', Gc2TribuneNavigator);