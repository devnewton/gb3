class Gc2Posts extends HTMLElement {

    constructor() {
        super();
    }

    insertPost(post) {
        let postElement = new Gc2Post();
        postElement.setup(post);
        this.appendChild(postElement);
        this.sortPosts();
    }

    sortPosts() {
        let postElements = this.getElementsByTagName('gc2-post');
        Array.prototype.slice.call(postElements).sort((a, b) => {
            if (a.time < b.time) {
                return 1;
            } else if (a.time > b.time) {
                return -1;
            } else if (a.tribune < b.tribune) {
                return 1;
            } else if (a.tribune > b.tribune) {
                return -1;
            } else if (a.id < b.id) {
                return 1;
            } else if(a.id > b.id) {
                return -1;
            }
            return 0;
        }).forEach((p) => {
            this.appendChild(p);
        });
    }
    insertPostAfterElement(post, element) {
        let postElement = new Gc2Post();
        postElement.setup(post);
        this.insertBefore(postElement, element.nextSibling);
    }

}
customElements.define('gc2-posts', Gc2Posts);