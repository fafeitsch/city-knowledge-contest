import {expect, Page} from '@playwright/test';
import selectors from './selectors';

export async function countDowns(users: { page: Page }[]) {
  // This is flaky
  // const countdownPromises = []
  // for (const {page} of users) {
  //   countdownPromises.push(expect(page.getByTestId('countdown-overlay')).toContainText('3'))
  //   countdownPromises.push(expect(page.getByTestId('countdown-overlay')).toContainText('2'))
  //   countdownPromises.push(expect(page.getByTestId('countdown-overlay')).toContainText('1'))
  // }
  // await Promise.all(countdownPromises)
  for (const {page} of users) {
    await expect(page.getByTestId(selectors.countDownOverlay)).toHaveCount(0)
    await expect(page.getByTestId(selectors.countDownOverlay)).toHaveCount(0)
  }
}

