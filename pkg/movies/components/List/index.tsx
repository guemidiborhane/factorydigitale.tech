
import { APISchemas } from '~/api'
import styles from './index.module.scss'
import { truncate } from '~/helpers'
import { fetchApi } from '~/helpers/http'
import Button from 'ui/Button'
import { useCan } from '~/hooks/useCan'
type Movie = APISchemas["github_com_guemidiborhane_factorydigitale_tech_pkg_movies_models.Movie"]
type Props = {
  movies: Movie[]
}
export default function ListMovies({ movies }: Props) {
  const canFavourite = useCan("movies:favourite")
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
    <div className="grid grid-cols-5 gap-3">
      {movies.map(movie => {
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
              {canFavourite && <Button onClick={favourite} data-movie-id={movie.id}>{movie.in_favourites ? '💔' : '❤️'}</Button>}
            </div>
          </div>
        )
      })}
    </div>
  )
}
