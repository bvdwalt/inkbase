import { test, expect } from "@playwright/test";
import { createPage, save } from "./helpers";

test.beforeEach(async ({ page }) => {
  // App-native confirm() dialogs (unsaved-changes guard, delete confirmation).
  page.on("dialog", (dialog) => dialog.accept());
  await page.goto("/");
});

test("[[ inserts a page link that navigates to the target page on click", async ({ page }) => {
  await createPage(page, "Page One");
  await save(page);

  await createPage(page, "Page Two");
  const editor = page.locator(".ProseMirror");
  await editor.click();
  await page.keyboard.type("[[Page One");
  await page.waitForSelector(".page-link-suggestion");
  await page.keyboard.press("Enter");
  await expect(page.locator(".ProseMirror a.page-link")).toHaveText("Page One");
  await save(page);

  await page.locator(".ProseMirror a.page-link").click();
  await expect(page.locator(".title-input")).toHaveValue("Page One");
});

test("typed markdown link syntax becomes a clickable external link opening a new tab", async ({
  page,
}) => {
  await createPage(page, "External Link Page");

  const editor = page.locator(".ProseMirror");
  await editor.click();
  await page.keyboard.type("[Example](https://example.com)");
  const link = page.locator('.ProseMirror a[href="https://example.com"]');
  await expect(link).toHaveText("Example");

  const [popup] = await Promise.all([page.waitForEvent("popup"), link.click()]);
  await popup.waitForLoadState().catch(() => {});
  expect(popup.url()).toBe("https://example.com/");
  await popup.close();
});

test("clicking a link to a deleted page shows an inline error instead of crashing", async ({
  page,
}) => {
  await createPage(page, "Soon Deleted");
  await save(page);

  await createPage(page, "Links To Deleted");
  const editor = page.locator(".ProseMirror");
  await editor.click();
  await page.keyboard.type("[[Soon Deleted");
  await page.waitForSelector(".page-link-suggestion");
  await page.keyboard.press("Enter");
  await save(page);

  const pageLink = page.locator(".ProseMirror a.page-link").first();
  const pageId = await pageLink.getAttribute("data-id");
  expect(pageId).toBeTruthy();

  // Simulate the linked page having been deleted server-side.
  await page.route(`**/api/pages/${pageId}`, (route) =>
    route.fulfill({ status: 404, contentType: "application/json", body: '{"error":"not found"}' }),
  );

  await pageLink.click();
  await expect(page.locator(".error")).toBeVisible();
});
