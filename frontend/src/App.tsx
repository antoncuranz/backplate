import {useRef, useState} from 'react'
import './App.css'
import {Button} from "@/components/ui/button.tsx";
import {Slider} from "@/components/ui/slider.tsx";
import {ToggleGroup, ToggleGroupItem} from "@/components/ui/toggle-group.tsx";
import Canvas from "@/components/Canvas.tsx";
import {Card, CardContent, CardFooter, CardHeader, CardTitle} from "@/components/ui/card.tsx";
import {toast, Toaster} from "sonner";

function App() {
  const colors = ["#000000", "#333333", "#666666", "#999999", "#CCCCCC", "#E6E6E6", "#F2F2F2", "#FFFFFF"];
  const initialLineWidth = 15;

  const [lineWidth, setLineWidth] = useState(initialLineWidth)
  const [color, setColor] = useState("black")
  const [scale, _] = useState(calculateScale)
  const canvasRef = useRef<HTMLCanvasElement>(null);

  function calculateScale() {
    if (window.innerWidth < 820)
      return window.innerWidth / 820;
    else
      return 1;
  }

  function clearImage() {
    const canvas = canvasRef.current!
    const context = canvas.getContext('2d')!
    context.fillStyle = 'white';
    context.fillRect(0, 0, canvas.width, canvas.height);
  }

  function submitImage() {
    toast("Hallo")
    canvasRef.current!.toBlob(function (blob) {
      const formData = new FormData();
      formData.append('image', blob!);

      fetch('http://192.168.1.20:8090/upload', {
        method: 'POST', body: formData,
      }).then(() => {
        clearImage()
        toast("Gemälde erfolgreich übermittelt!")
      }).catch(toast);
    })
  }

  return (<>
    <Card>
      <CardHeader>
        <CardTitle>Hey! 👋🏻</CardTitle>
      </CardHeader>
      <CardContent>
        <form>
          <div className="grid w-full items-center gap-4">
            <div className="flex flex-col space-y-1.5">
              <ToggleGroup variant="outline" type="single" onValueChange={x => setColor(x)}>
                {colors.map(color => (
                  <ToggleGroupItem key={color} value={color} className={"color-button"}>
                    <div className="color-div" style={{background: color, borderColor: color}}></div>
                  </ToggleGroupItem>
                ))}
              </ToggleGroup>
              <Slider className="slider" defaultValue={[initialLineWidth]} min={1} max={72}
                      onValueChange={x => setLineWidth(x[0])}/>
            </div>
          </div>
        </form>
      </CardContent>
      <CardFooter className="flex justify-between">
        <div>
          <Button variant="outline" style={{marginRight: "10px"}} onClick={clearImage}>Clear</Button>
          {/*<Button variant="outline">Undo</Button>*/}
        </div>
        <Button onClick={submitImage}>Submit</Button>
      </CardFooter>
    </Card>
    <Canvas ref={canvasRef} lineWidth={lineWidth} color={color} scale={scale}></Canvas>
    <Toaster />
  </>)
}

export default App