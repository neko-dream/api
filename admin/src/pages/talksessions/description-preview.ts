import xss from "xss";

// web のセッション画面（app/components/features/xss-html）と同じ許可リスト。
// 許可されないタグは web と同じく文字として表示されるので、プレビューも同じ見え方になる。
const WEB_ALLOW_LIST = {
  img: ["style", "src", "alt"],
  p: [],
  a: ["href", "target", "rel"],
  span: [],
};

export function sanitizeDescription(html: string): string {
  return xss(html, { whiteList: WEB_ALLOW_LIST });
}

// サニタイズに加えて、sandbox 付き iframe の中で CSP をかけて表示する（スクリプトは実行させない）
const CSP = "default-src 'none'; img-src https: data:; style-src 'unsafe-inline'";
const STYLE = `
  body { margin: 16px; font-family: -apple-system, BlinkMacSystemFont, "Hiragino Sans", sans-serif; line-height: 1.7; color: #191c1e; }
  img { max-width: 100%; height: auto; }
  a { color: #007aff; text-decoration: underline; }
  p { margin: 0 0 1em; }
`;

export function buildPreviewDocument(html: string): string {
  return `<!doctype html><html lang="ja"><head><meta charset="utf-8"><meta http-equiv="Content-Security-Policy" content="${CSP}"><base target="_blank"><style>${STYLE}</style></head><body>${sanitizeDescription(html)}</body></html>`;
}
