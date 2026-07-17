import type { Page } from "@playwright/test";

export async function createPage(page: Page, title: string): Promise<void> {
  await page.getByRole("button", { name: "+ New page" }).click();
  await page.waitForSelector(".title-input");
  await page.locator(".title-input").fill(title);
}

export async function save(page: Page): Promise<void> {
  await Promise.all([
    page.waitForResponse(
      (r) => r.request().method() === "PUT" && r.url().includes("/api/pages/"),
    ),
    page.locator(".btn-primary").click(),
  ]);
}
