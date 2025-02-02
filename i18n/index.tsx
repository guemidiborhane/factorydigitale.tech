import { useT as useTr, type Autocomplete, type TParams, tr, Talkr } from "talkr";
import type { ReactNode } from "preact/compat";
import languages, { defaultLanguage } from './langs'

const dl = languages[defaultLanguage]
export type Key = Autocomplete<typeof dl> | string;

export const useT = () => {
  const { locale, setLocale, languages, defaultLanguage } = useTr();
  return {
    setLocale,
    locale,
    t: (key: Key, params?: TParams) =>
      tr({ locale, languages, defaultLanguage }, key, { count: params?.other && 11, ...params }) || key.split('.').slice(-1),
  };
};

export default function I18n({ children }: { children: ReactNode }) {
  return (
    <Talkr languages={languages} defaultLanguage={defaultLanguage}>
      {children}
    </Talkr>
  )
}
