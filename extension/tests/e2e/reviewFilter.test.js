describe('Review Filter Extension', () => {
    beforeAll(async () => {
        await page.goto('https://www.glassdoor.com/Reviews/test-company');
    });

    test('hides fake reviews and shows flash message', async () => {
        const reviewDiv = await page.$('[data-test="review-details-container"]');
        const isHidden = await reviewDiv.evaluate(el => el.style.display === 'none');
        expect(isHidden).toBe(true);

        const flashMessage = await page.$('.flash-message');
        expect(flashMessage).toBeTruthy();
    });

    test('shows review when clicking Show Review', async () => {
        await page.click('[data-testid="show-review"]');
        const reviewDiv = await page.$('[data-test="review-details-container"]');
        const isVisible = await reviewDiv.evaluate(el => el.style.display === 'block');
        expect(isVisible).toBe(true);
    });
});