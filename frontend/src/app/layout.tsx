import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Astro-Karuta - 天文かるた",
  description: "天文知識をかるた形式で学ぶ子ども向け天文教育ゲーム",
  manifest: "/manifest.json",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="ja">
      <body>{children}</body>
    </html>
  );
}
