import assert from "node:assert/strict";
import test from "node:test";
import {
  buildPreviewDocument,
  sanitizeDescription,
} from "../src/pages/talksessions/description-preview.ts";

test("description is sanitized with the same allow list as the web session page", () => {
  const html = sanitizeDescription(
    '<p>本文<span>強調</span></p><img src="https://example.com/a.png" alt="画像" style="width:10px" onerror="alert(1)"><a href="https://example.com" target="_blank" rel="noopener" onclick="alert(1)">リンク</a>',
  );
  assert.match(html, /<p>本文<span>強調<\/span><\/p>/);
  assert.match(html, /<img src="https:\/\/example.com\/a.png" alt="画像" style="width:10px;">/);
  assert.match(html, /<a href="https:\/\/example.com" target="_blank" rel="noopener">リンク<\/a>/);
  assert.doesNotMatch(html, /onerror|onclick/);
});

test("tags the web does not allow are shown as text, as on the web", () => {
  const html = sanitizeDescription("<h2>見出し</h2><script>alert(1)</script>");
  assert.doesNotMatch(html, /<h2>|<script>/);
  assert.match(html, /&lt;h2&gt;見出し&lt;\/h2&gt;/);
});

test("preview document blocks scripts and only loads https images", () => {
  const doc = buildPreviewDocument("<p>本文</p>");
  assert.match(doc, /<p>本文<\/p>/);
  const csp = doc.match(/http-equiv="Content-Security-Policy" content="([^"]+)"/)?.[1] ?? "";
  assert.match(csp, /default-src 'none'/);
  assert.match(csp, /img-src https: data:/);
  assert.doesNotMatch(csp, /script-src/);
  assert.match(doc, /<base target="_blank">/);
});
