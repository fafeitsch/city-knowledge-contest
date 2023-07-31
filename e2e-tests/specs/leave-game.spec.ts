import {Browser, expect, test} from '@playwright/test';
import {countDowns} from './utils';
import selectors from './selectors';

async function createRoom(browser: Browser) {
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

test('should be able to leave the game while in waiting room', async ({browser}) => {
  const {alice, bob} = await createRoom(browser);

  await bob.getByTestId(selectors.leaveGameButton).click()

  await expect(alice.getByTestId(selectors.playerListEntry)).toHaveCount(1)
  await expect(bob.getByTestId(selectors.userNameInput)).toHaveValue('Bob')
})

test('it should be able to leave the game while in game mode', async ({browser}) => {
  const {alice, bob} = await createRoom(browser);
  await alice.getByTestId(selectors.startGameButton).click()

  await countDowns([{page: alice}, {page: bob}])

  await alice.mouse.click(400,400)
  await expect(alice.getByTestId(selectors.playerListEntry).nth(0)).toHaveText(/Alice 0 Punkte \+ \d\d/)

  await bob.getByTestId(selectors.leaveGameButton).click()

  await expect(alice.getByTestId(selectors.playerListEntry)).toHaveCount(1)
  await expect(bob.getByTestId(selectors.userNameInput)).toHaveValue('Bob')
})

test('it should be able to join a game late', async ({browser}) => {
  const {alice, bob} = await createRoom(browser);
  await bob.getByTestId(selectors.startGameButton).click()
  await countDowns([{page: alice}, {page: bob}])
  await alice.mouse.click(400,400)
  await expect(alice.getByTestId(selectors.playerListEntry).nth(0)).toHaveText(/Alice 0 Punkte \+ \d\d/)
  await bob.mouse.click(400,400)
  await expect(bob.getByTestId(selectors.playerListEntry).nth(1)).toHaveText(/Bob \d\d Punkte/)
  await alice.getByTestId(selectors.proceedGameButton).click()
  await countDowns([{page: alice}, {page: bob}])

  const contextCarol = await browser.newContext();
  const carol = await contextCarol.newPage();

  await carol.goto(alice.url());
  await carol.getByTestId(selectors.userNameInput).fill('Carol');
  await carol.getByTestId(selectors.joinRoomButton).click()

  await expect(carol.getByTestId(selectors.questionCard)).toContainText('Suche den Ort Frankfurter Straße')
  await expect(carol.getByTestId(selectors.playerListEntry).nth(0)).toHaveText(/Alice 0 Punkte/)
  await expect(carol.getByTestId(selectors.playerListEntry).nth(1)).toHaveText(/Bob 0 Punkte/)
  await expect(carol.getByTestId(selectors.playerListEntry).nth(2)).toHaveText(/Carol 0 Punkte/)
  await expect(carol.getByTestId(selectors.playerListEntry).nth(2)).toHaveClass(/highlight/)

  // Points are only updated at the end of every round, so we let all players click
  await carol.mouse.click(500,300)
  await bob.mouse.click(500,300)
  await alice.mouse.click(500,300)

  await expect(carol.getByTestId(selectors.playerListEntry).nth(0)).toHaveText(/Alice \d\d Punkte/)
  await expect(carol.getByTestId(selectors.playerListEntry).nth(1)).toHaveText(/Bob \d\d Punkte/)
  await expect(carol.getByTestId(selectors.playerListEntry).nth(2)).toHaveText(/Carol 0 Punkte/)
  await expect(carol.getByTestId(selectors.playerListEntry).nth(2)).toHaveClass(/highlight/)

  await expect(bob.getByTestId(selectors.playerListEntry).nth(0)).toHaveText(/Alice \d\d Punkte/)
  await expect(bob.getByTestId(selectors.playerListEntry).nth(1)).toHaveText(/Bob \d\d Punkte/)
  await expect(bob.getByTestId(selectors.playerListEntry).nth(1)).toHaveClass(/highlight/)
  await expect(bob.getByTestId(selectors.playerListEntry).nth(2)).toHaveText(/Carol 0 Punkte/)
})
