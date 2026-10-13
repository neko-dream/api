// Match Go strings.TrimSpace and RuneCountInString, not HTML's UTF-16 maxLength.
const edgeWhitespace = /^\p{White_Space}+|\p{White_Space}+$/gu;
const passwordLoginEmail = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;

export function trimUserField(value: string): string {
  return value.replace(edgeWhitespace, "");
}

export function validateUserDisplayName(value: unknown): string {
  const name = typeof value === "string" ? trimUserField(value) : "";
  if (!name) return "表示名を入力してください";
  return Array.from(name).length > 50
    ? "表示名は50文字以内で入力してください"
    : "";
}

export function validateUserEmail(value: unknown): string {
  const email = typeof value === "string" ? trimUserField(value) : "";
  // These ASCII dot-atom restrictions match the server's net/mail + login format.
  const parts = email.split("@");
  if (
    !passwordLoginEmail.test(email) ||
    parts.some(
      (part) =>
        part.startsWith(".") || part.endsWith(".") || part.includes(".."),
    )
  )
    return "メールアドレスを入力してください";
  return email.length > 254
    ? "メールアドレスは254文字以内で入力してください"
    : "";
}
