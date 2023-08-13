import {expect, test} from '@playwright/test';
import {countDowns, createRoom} from './utils';
import selectors from './selectors';

test('should timeout kick button if user waits too long', async ({browser}) => {
  const {alice, bob} = await createRoom(browser);
  await alice.getByTestId(selectors.kickPlayerButton).click()
  await expect(alice.getByTestId(selectors.kickPlayerButton)).toContainText('Kick?')
  await expect(alice.getByTestId(selectors.kickPlayerButton)).toContainText('Kick!')
  await expect(alice.getByTestId(selectors.kickPlayerButton)).toContainText('🗙')
  await alice.getByTestId(selectors.kickPlayerButton).click()
  await expect(alice.getByTestId(selectors.kickPlayerButton)).toContainText('Kick?')
  await expect(alice.getByTestId(selectors.kickPlayerButton)).toContainText('Kick!')
})

test('should be possible to kick player in waiting room', async ({browser}) => {
  const {alice, bob} = await createRoom(browser);

  await expect(alice.getByTestId(selectors.kickPlayerButton)).toHaveCount(1);
  await alice.getByTestId(selectors.kickPlayerButton).click()
  await expect(alice.getByTestId(selectors.kickPlayerButton)).toContainText('Kick?')
  await expect(alice.getByTestId(selectors.kickPlayerButton)).toContainText('Kick!')
  await alice.getByTestId(selectors.kickPlayerButton).click()

  await expect(alice.getByTestId(selectors.playerListEntry)).toHaveCount(1)

  await expect(bob.getByTestId(selectors.userNameInput)).toBeVisible()
  await expect(bob.getByTestId(selectors.userNameInput)).toHaveValue('Bob')
  await expect(bob.url()).not.toContain('/room')
})

test('should be possible to kick player in running game', async ({browser}) => {
  const {alice, bob} = await createRoom(browser);

  await bob.getByTestId(selectors.startGameButton).click()

  await countDowns([{page: alice}, {page: bob}])

  await expect(alice.getByTestId(selectors.kickPlayerButton)).toHaveCount(1);
  await alice.getByTestId(selectors.kickPlayerButton).click()
  await expect(alice.getByTestId(selectors.kickPlayerButton)).toContainText('Kick?')
  await expect(alice.getByTestId(selectors.kickPlayerButton)).toContainText('Kick!')
  await alice.getByTestId(selectors.kickPlayerButton).click()

  await expect(alice.getByTestId(selectors.playerListEntry)).toHaveCount(1)

  await expect(bob.getByTestId(selectors.userNameInput)).toBeVisible()
  await expect(bob.getByTestId(selectors.userNameInput)).toHaveValue('Bob')
  await expect(bob.url()).not.toContain('/room')
})
