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
        
        let tribuneSpan = document.createElement("span");
        tribuneSpan.innerText = tribuneName;
        tribuneButton.appendChild(tribuneSpan);

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

    /**
     * 
     * @param {string} tribuneName 
     * @param {boolean} hasBigorno 
     * @param {boolean} hasReply 
     */
    setTribuneNotifications(tribuneName, hasBigorno, hasReply) {
        this.querySelectorAll("gc2-tribune-navigator-button").forEach((element) => {
            if(element.tribuneName === tribuneName) {
                element.toggleAttribute("bigorno", hasBigorno);
                element.toggleAttribute("reply", hasReply);
            }
        });
    }
}

customElements.define('gc2-tribune-navigator', Gc2TribuneNavigator);