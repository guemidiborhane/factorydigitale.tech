import { useT } from "i18n"

type Direction = "rtl" | "ltr"

export default function useDirection(): Direction {
  const { locale } = useT()
  return locale === "ar" ? "rtl" : "ltr"
}

