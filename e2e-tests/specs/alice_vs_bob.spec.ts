import {test, expect, Page} from '@playwright/test';
import {countDowns} from './utils';
import {default as selectors} from './selectors';

test('should let alice play versus bob two rounds', async ({browser}) => {
  const contextAlice = await browser.newContext();
  const alice = await contextAlice.newPage();

  const contextBob = await browser.newContext();
  const bob = await contextBob.newPage();

  const users = [{name: 'Alice', page: alice}, {name: 'Bob', page: bob}]

  await alice.goto('http://localhost:5173');
  await alice.getByTestId(selectors.userNameInput).fill('Alice');
  expect(await alice.getByTestId(selectors.joinRoomButton).count()).toEqual(0)
  await alice.getByTestId(selectors.createRoomButton).click();

  await expect(alice.getByTestId('room-link-display')).toHaveText(/^http:\/\/localhost:5173\/room/)
  const link = await (alice.getByTestId(selectors.roomLinkDisplay).textContent())

  await bob.goto(link)

  await bob.getByTestId(selectors.userNameInput).fill('Bob');
  expect(await bob.getByTestId(selectors.createRoomButton).count()).toEqual(0)
  await bob.getByTestId(selectors.joinRoomButton).click();

  for (const user of users) {
    const playerList = await user.page.getByTestId(selectors.playerListEntry)
    expect(await playerList.count()).toEqual(2)
    await expect(playerList.nth(0)).toContainText('Alice')
    await expect(playerList.nth(1)).toContainText('Bob')
    const startButton = await user.page.getByTestId(selectors.startGameButton)
    await expect(startButton).toBeDisabled()
  }

  await bob.getByTestId('select-streetlist').selectOption({label: 'Würzburg Altstadt'})

  await expect(alice.getByTestId(selectors.selectStreetList)).toHaveValue('wuerzburg-altstadt.json')

  await expect(bob.getByTestId(selectors.startGameButton)).toBeEnabled()
  await expect(alice.getByTestId(selectors.startGameButton)).toBeEnabled()

  await alice.getByTestId(selectors.numberOfQuestionsInput).clear()
  await alice.getByTestId(selectors.numberOfQuestionsInput).type('2')

  await expect(bob.getByTestId(selectors.numberOfQuestionsInput)).toHaveValue('2')

  await bob.getByTestId(selectors.maxAnswerTimeInput).clear()
  await bob.getByTestId(selectors.maxAnswerTimeInput).type('10')

  await expect(alice.getByTestId(selectors.maxAnswerTimeInput)).toHaveValue('10')

  await bob.getByTestId(selectors.startGameButton).click();
  await countDowns(users);

  //dummy click since drag and drop doesn't work in playwright yet
  await alice.mouse.click(400, 300)
  await expect(alice.getByTestId(selectors.playerListEntry).nth(0)).toHaveText(/0 Punkte + 5|6|7|8|9\d/)
  await expect(alice.getByTestId(selectors.gameStateCard)).toHaveText('Richtig')
  await expect(alice.getByTestId(selectors.proceedGameButton)).toHaveCount(0)

  await bob.waitForTimeout(1000)
  await bob.mouse.click(400, 300)
  await expect(bob.getByTestId(selectors.playerListEntry).nth(1)).toHaveText(/0 Punkte + 3|4|5|6\d/)
  await expect(bob.getByTestId(selectors.gameStateCard)).toContainText('Richtig')
  await expect(bob.getByTestId(selectors.proceedGameButton)).toHaveCount(1)
  await expect(alice.getByTestId(selectors.proceedGameButton)).toHaveCount(1)

  await alice.getByTestId(selectors.proceedGameButton).click()
  await countDowns(users);

  await expect(alice.getByTestId(selectors.questionCard)).toHaveText('Suche den Ort Frankfurter Straße')
  await expect(bob.getByTestId(selectors.questionCard)).toHaveText('Suche den Ort Frankfurter Straße')

  //dummy click since drag and drop doesn't work in playwright yet
  await alice.mouse.click(400, 300)
  await expect(alice.getByTestId(selectors.playerListEntry).nth(0)).toHaveText(/5|6|7|8|9\d Punkte + 0/)
  await expect(alice.getByTestId(selectors.gameStateCard)).toHaveText('Falsch')
  await expect(alice.getByTestId(selectors.proceedGameButton)).toHaveCount(0)

  await bob.waitForTimeout(1000)
  await bob.mouse.click(400, 300)
  await expect(bob.getByTestId(selectors.playerListEntry).nth(1)).toHaveText(/3|4|5|6\d Punkte + 0/)
  await expect(bob.getByTestId(selectors.gameStateCard)).toContainText('Falsch')
  await expect(bob.getByTestId(selectors.proceedGameButton)).toHaveCount(1)
  await expect(alice.getByTestId(selectors.proceedGameButton)).toHaveCount(1)

  await bob.getByTestId('proceed-game-button').click()

  await expect(bob.getByTestId(selectors.gameOverTitle)).toHaveText('Das Spiel ist leider vorbei.')
  await expect(alice.getByTestId(selectors.gameOverTitle)).toHaveText('Das Spiel ist leider vorbei.')

  await expect(bob.getByTestId(selectors.playerListEntry).nth(0)).toHaveText(/Alice\s+6|7|8|9\d Punkte/)
  await expect(bob.getByTestId(selectors.playerListEntry).nth(1)).toHaveText(/Bob\s+3|4|5|6\d Punkte/)
  await expect(bob.getByTestId(selectors.playerListEntry).nth(0)).not.toHaveClass(/highlight/)
  await expect(bob.getByTestId(selectors.playerListEntry).nth(1)).toHaveClass(/highlight/)

  await expect(alice.getByTestId(selectors.playerListEntry).nth(0)).toHaveText(/Alice\s+6|7|8|9\d Punkte/)
  await expect(alice.getByTestId(selectors.playerListEntry).nth(1)).toHaveText(/Bob\s+3|4|5|6\d Punkte/)
  await expect(alice.getByTestId(selectors.playerListEntry).nth(0)).toHaveClass(/highlight/)
  await expect(alice.getByTestId(selectors.playerListEntry).nth(1)).not.toHaveClass(/highlight/)

  await alice.getByTestId(selectors.startNewGameButton).click()
  await bob.getByTestId(selectors.startNewGameButton).click()

  await expect(alice.getByTestId(selectors.userNameInput)).toHaveValue('Alice')
  await expect(bob.getByTestId(selectors.userNameInput)).toHaveValue('Bob')
});
