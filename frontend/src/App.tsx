import {useRef, useState} from 'react'
import './App.css'
import {Button} from "@/components/ui/button.tsx";
import {Slider} from "@/components/ui/slider.tsx";
import {ToggleGroup, ToggleGroupItem} from "@/components/ui/toggle-group.tsx";
import Canvas from "@/components/Canvas.tsx";
import {Card, CardContent} from "@/components/ui/card.tsx";
import {toast, Toaster} from "sonner";
import {Separator} from "@/components/ui/separator.tsx";
import {Label} from "@/components/ui/label.tsx";

function App() {
  const colors = ["#000000", "#333333", "#666666", "#999999", "#CCCCCC", "#E6E6E6", "#F2F2F2", "#FFFFFF"];
  const initialLineWidth = 40;

  const [lineWidth, setLineWidth] = useState(initialLineWidth)
  const [color, setColor] = useState("black")
  const canvasRef = useRef<HTMLCanvasElement>(null);

  function clearImage() {
    const canvas = canvasRef.current!
    const context = canvas.getContext('2d')!
    context.fillStyle = 'white';
    context.fillRect(0, 0, canvas.width, canvas.height);
  }

  function submitImage() {
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
    <div className="hidden h-full flex-col md:flex">
      <div
        className="container flex flex-col items-start justify-between space-y-2 py-4 sm:flex-row sm:items-center sm:space-y-0 md:h-16">
        <h2 className="text-lg font-semibold">Backplate</h2>
        <div className="ml-auto flex w-full space-x-2 sm:justify-end">
          preset was here
        </div>
      </div>
      <Separator/>
    </div>
    <div className="container gap-4 main-container">
      <Card className="canvas-card">
        <CardContent style={{padding: 0, height: "100%"}}>
          <Canvas ref={canvasRef} lineWidth={lineWidth} color={color}></Canvas>
        </CardContent>
      </Card>
      <div className="pt-4 tools-container">
        <div className="grid gap-4" style={{marginBottom: "14px"}}>
          <div className="flex items-center justify-between">
            <Label>Color</Label>
          </div>
          <ToggleGroup id="colors-togglegroup" variant="outline" type="single" onValueChange={x => setColor(x)} defaultValue="#000000">
            {colors.map(color => (
              <ToggleGroupItem key={color} value={color} className={"color-button"}>
                <div className="color-div" style={{background: color, borderColor: color}}></div>
              </ToggleGroupItem>
            ))}
          </ToggleGroup>
        </div>

        <div className="grid gap-4" style={{marginBottom: "28px"}}>
          <div className="flex items-center justify-between">
            <Label htmlFor="lineWidth">Line Width</Label>
            <span
              className="w-12 rounded-md border border-transparent px-2 py-0.5 text-right text-sm text-muted-foreground hover:border-border">
                {lineWidth}
              </span>
          </div>
          <Slider id="lineWidth" className="slider" defaultValue={[initialLineWidth]} min={20} max={200}
                  onValueChange={x => setLineWidth(x[0])}/>
        </div>
        <div className="gap-4 button-container" style={{display: "flex"}}>
          <Button variant="outline" onClick={clearImage}>Clear</Button>
          {/*<Button variant="outline">Undo</Button>*/}
          <Button onClick={submitImage}>Submit</Button>
        </div>
      </div>
    </div>
    <Toaster/>
  </>)
}

export default App