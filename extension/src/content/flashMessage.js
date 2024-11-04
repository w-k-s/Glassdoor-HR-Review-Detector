export class FlashMessage {
    constructor(onShowReview, onMarkGenuine, onMarkFake) {
        this.onShowReview = onShowReview;
        this.onMarkGenuine = onMarkGenuine;
        this.onMarkFake = onMarkFake;
    }

    create() {
        const div = document.createElement('div');
        div.className = 'flash-message';
        div.innerHTML = this.getHiddenTemplate();
        this.attachHiddenListeners(div);
        return div;
    }

    showReviewMode(flashDiv) {
        flashDiv.innerHTML = this.getShowingTemplate();
        this.attachShowingListeners(flashDiv);
    }

    reset(flashDiv) {
        flashDiv.innerHTML = this.getHiddenTemplate();
        this.attachHiddenListeners(flashDiv);
    }

    private getHiddenTemplate() {
        return `
        <a class="show-review" data-testid="show-review">Show Review</a>
        <span>This review has been hidden because it seems to be fake.</span>
      `;
    }

    private getShowingTemplate() {
        return `
        <a class="genuine-review" data-testid="mark-genuine">Seems Genuine</a>
        <a class="fake-review" data-testid="mark-fake">Seems Fake</a>
        <span>This review has been hidden because it seems to be fake.</span>
      `;
    }

    private attachHiddenListeners(div) {
        div.querySelector('.show-review').addEventListener('click', this.onShowReview);
    }

    private attachShowingListeners(div) {
        div.querySelector('.genuine-review').addEventListener('click', this.onMarkGenuine);
        div.querySelector('.fake-review').addEventListener('click', this.onMarkFake);
    }
}