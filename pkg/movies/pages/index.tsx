import { type LoaderFunctionArgs, redirect } from "react-router-dom";
import { useLoaderData } from "~/hooks";
import { fetchApi } from "~/helpers/http";
import { useCan } from "~/hooks/useCan";
import { useT } from "i18n/index";
import styles from './Index.module.scss'
import { truncate } from "~/helpers";
import Button from "ui/Button";
import { useComputed, useSignal } from "@preact/signals";

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
  const canFavourite = useCan("movies:favourite")
  const { t } = useT()

  const loadMoreMovies = () => {
    fetchApi('/api/movies', { query: { offset: offset.value } }).then(response => {
      if (response.ok) {
        movies.value = [...movies.value, ...response.data]
      }
    })
  }

  const favourite = (event: MouseEvent) => {
    const target = (event.target as HTMLButtonElement)
    const movie_id = parseInt(target.dataset['movieId'] as string)
    fetchApi('/api/movies/favourite', { method: 'post', body: { movie_id } }).then(response => {
      if (response.ok) {
        const { data } = response
        console.log(data)

        target.innerHTML = data.id == 0 ? '❤️' : '💔'
      }
    })
  }

  return (
    <>
      <div className="grid grid-cols-5 gap-3">
        {movies.value.map(movie => {
          return (
            <div className="box !p-0 !border-gray-800 overflow-y-auto" style={{
              background: `url(${movie.poster})`,
              backgroundSize: 'contain'
            }}>
              <div className={styles.Shadow}>
                <h3 className="font-sans text-white font-semibold text-2xl">{movie.title}</h3>
                <p>{movie.overview && truncate(movie.overview, 70)}</p>
                <div className="flex flex-row gap-3 flex-wrap justify-center">
                  {movie.genres && movie.genres.map(genre => <span className="bg-yellow-300 rounded-full text-xs text-black py-1 px-2 font-semibold whitespace-nowrap">{genre}</span>)}
                </div>
                <Button onClick={favourite} data-movie-id={movie.id}>{movie.in_favourites ? '💔' : '❤️'}</Button>
              </div>
            </div>
          )
        })}

      </div>
      <div className="flex justify-center">
        <Button onClick={loadMoreMovies}>{t('misc.load_more')}</Button>
      </div>
    </>
  )
}
