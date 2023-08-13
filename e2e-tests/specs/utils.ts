import {Browser, expect, Page} from '@playwright/test';
import selectors from './selectors';

export async function createRoom(browser: Browser) {
  const contextAlice = await browser.newContext();
  const alice = await contextAlice.newPage();

  const contextBob = await browser.newContext();
  const bob = await contextBob.newPage();

  await alice.goto('http://localhost:5173');
  await alice.getByTestId(selectors.userNameInput).fill('Alice');
  await alice.getByTestId(selectors.createRoomButton).click();

  await expect(alice.getByTestId(selectors.roomLinkDisplay)).toHaveText(/^http:\/\/localhost:5173\/room/)
  const link = await (alice.getByTestId(selectors.roomLinkDisplay).textContent())

  await bob.goto(link)

  await bob.getByTestId(selectors.userNameInput).fill('Bob');
  await bob.getByTestId(selectors.joinRoomButton).click();

  await bob.getByTestId(selectors.selectStreetList).selectOption({label: 'Würzburg Altstadt'})
  return {alice, bob};
}

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

