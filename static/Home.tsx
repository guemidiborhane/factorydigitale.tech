import RootLayout from "./RootLayout";
import { useWebSocket } from "ws/client";


export default function Home() {
  useWebSocket<string>({
    channel: "news",
    receiver(data) {
      console.log(data, "news")
    },
  })

  useWebSocket<string>({
    channel: "test",
    receiver(data) {
      console.log(data, "testHome")
    }
  })

  return (
    <RootLayout>
      <div class="box">
        HOME !!!
      </div>
    </RootLayout>
  )
}
