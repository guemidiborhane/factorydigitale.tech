import { Link, type LoaderFunctionArgs, redirect } from "react-router-dom";
import { useLoaderData } from "~/hooks";
import { fetchApi } from "~/helpers/http";
import { useT } from "i18n/index";
import Button from "ui/Button";
import { useComputed, useSignal } from "@preact/signals";
import { ListMovies } from "@/movies/components";
import { routes } from "~/router";

export async function loader({ request }: LoaderFunctionArgs) {
  const response = await fetchApi('/api/movies', { signal: request.signal })
  if (!response.ok) return redirect('/')
  return response
}

export default function IndexMovies() {
  const response = useLoaderData<'/api/movies'>()

  if (!response.ok) return null

  const movies = useSignal(response.data)
  const offset = useComputed(() => movies.value.length)
  const { t } = useT()

  const loadMoreMovies = () => {
    fetchApi('/api/movies', { query: { offset: offset.value } }).then(response => {
      if (response.ok) {
        movies.value = [...movies.value, ...response.data]
      }
    })
  }

  return (
    <>
      <ListMovies movies={movies.value} />
      <div className="flex justify-center">
        <Button onClick={loadMoreMovies}>{t('misc.load_more')}</Button>
      </div>
    </>
  )
}
