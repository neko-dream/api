import assert from "node:assert/strict";
import test from "node:test";
import {
  trimUserField,
  validateUserDisplayName,
  validateUserEmail,
} from "../src/pages/users/validation.ts";

test("managed display names use the server's Unicode character count and edge trimming", () => {
  for (const value of [
    undefined,
    null,
    "",
    " ",
    "\u0085\u3000",
    "名".repeat(51),
    "😀".repeat(51),
  ]) {
    assert.notEqual(validateUserDisplayName(value), "");
  }
  for (const value of [
    " 表示名 ",
    "名".repeat(50),
    "😀".repeat(50),
    "\u0085表示名\u0085",
  ]) {
    assert.equal(validateUserDisplayName(value), "");
  }
  assert.equal(trimUserField("\u0085 表示名 \u0085"), "表示名");
});

test("managed email validation follows the password login ASCII format and bare dot-atoms", () => {
  for (const value of [
    undefined,
    "",
    "bad",
    "Name <user@example.com>",
    '"user"@example.com',
    "u@example.com\r\nBcc: x@example.com",
    ".user@example.com",
    "user.@example.com",
    "u..ser@example.com",
    "user@.example.com",
    "user@example..com",
    "user@localhost",
    "user@例え.test",
    "a".repeat(243) + "@example.com",
  ]) {
    assert.notEqual(validateUserEmail(value), "");
  }
  for (const value of [
    "user@example.com",
    "user+tag@example.com",
    " user@example.com ",
    "\u0085user@example.com\u0085",
    "user.name@sub.example.com",
    "USER@EXAMPLE.COM",
    "a".repeat(242) + "@example.com",
  ]) {
    assert.equal(validateUserEmail(value), "");
  }
});
