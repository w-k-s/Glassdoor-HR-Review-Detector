import { ReviewFilter } from '../../src/content/ReviewFilter';

describe('ReviewFilter', () => {
    let reviewFilter;

    beforeEach(() => {
        document.body.innerHTML = `
      <div data-test="review-details-container" id="empReview_123"></div>
    `;
        reviewFilter = new ReviewFilter();
    });

    test('finds review divs correctly', () => {
        const reviewDivs = reviewFilter.findReviewDivs();
        expect(reviewDivs).toHaveLength(1);
        expect(reviewDivs[0].id).toBe('empReview_123');
    });
});