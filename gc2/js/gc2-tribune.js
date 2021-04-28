import {Gc2Post} from "./gc2-post.js"

export class Gc2Tribune extends HTMLElement {

    /**
     * @type boolean
     */
    hasBigorno = false;

    /**
     * @type boolean
     */
    hasReply = false;

    lastReadPostId = -1;

    constructor() {
        super();
    }

    insertPost(post) {
        if (!document.getElementById(`${post.id}@${post.tribune}`)) {
            let postElement = new Gc2Post();
            postElement.setup(post);
            this.appendChild(postElement);
            this.sortPosts();
        }
    }

    sortPosts() {
        let postElements = this.getElementsByTagName('gc2-post');
        Array.prototype.slice.call(postElements).sort((a, b) => {
            if (a.postId < b.postId) {
                return -1;
            } else if (a.postId > b.postId) {
                return 1;
            }
            return 0;
        }).forEach((p) => {
            this.appendChild(p);
        });

        this.hasBigorno = false;
        this.hasReply = false;
        for (let i = postElements.length - 1; i > 0; --i) {
            let post = postElements[i];
            if(post.id < this.lastReadPostId || post.isMine) {
                break;
            }
            this.hasReply |= post.isReply;
            this.hasBigorno |= post.isBigorno;
        }
    }

    markAsRead() {
        this.hasBigorno = false;
        this.hasReply = false;
        let lastPost = this.querySelector('gc2-post:last-of-type');
        if(lastPost) {
            this.lastReadPostId = lastPost.id;
        }
    }
}
customElements.define('gc2-tribune', Gc2Tribune);

class Gc2Tribunes extends HTMLElement {
    constructor() {
        super();
    }
}
customElements.define('gc2-tribunes', Gc2Tribunes);