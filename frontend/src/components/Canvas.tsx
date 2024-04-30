import React, {useRef, useEffect, MouseEvent, TouchEvent, useImperativeHandle} from 'react'

interface Props {
  lineWidth: number,
  color: string,
  scale: number
}

const Canvas = React.forwardRef<HTMLCanvasElement, Props>(
  ({lineWidth, color, scale}: Props, forwardedRef: React.ForwardedRef<HTMLCanvasElement>) => {

    const canvasRef = useRef<HTMLCanvasElement>(null);
    useImperativeHandle(forwardedRef, () => canvasRef.current!);

    useEffect(() => {
      const canvas = canvasRef.current!
      const context = canvas.getContext('2d')!

      context.fillStyle = 'white';
      context.fillRect(0, 0, canvas.width, canvas.height);
      context.lineCap = 'round';

      // Prevent scrolling when touching the canvas
      document.body.addEventListener("touchstart", function (e) {
        if (e.target == canvas) e.preventDefault();
      }, { passive: false });
      document.body.addEventListener("touchend", function (e) {
        if (e.target == canvas) e.preventDefault();
      }, { passive: false });
      document.body.addEventListener("touchmove", function (e) {
        if (e.target == canvas) e.preventDefault();
      }, { passive: false });
    }, [])

    useEffect(() => {
      const context = canvasRef.current!.getContext('2d')!
      context.lineWidth = lineWidth
    }, [lineWidth]);

    useEffect(() => {
      const context = canvasRef.current!.getContext('2d')!
      context.strokeStyle = color;
    }, [color]);

    let x = 0, y = 0;
    let isMouseDown = false;

    function startDrawing(event: MouseEvent) {
      isMouseDown = true;
      x = (event.nativeEvent.offsetX) / scale;
      y = (event.nativeEvent.offsetY) / scale;
    }

    function startDrawingTouch(event: TouchEvent) {
      isMouseDown = true;
      const rect = canvasRef.current!.getBoundingClientRect()
      const touch = event.targetTouches[0];
      x = (touch.clientX - rect.x) / scale;
      y = (touch.clientY - rect.y) / scale;
    }

    function drawLine(event: MouseEvent) {
      if (isMouseDown) {
        const newX = event.nativeEvent.offsetX / scale;
        const newY = event.nativeEvent.offsetY / scale;
        const context = canvasRef.current!.getContext("2d")!
        context.beginPath();
        context.moveTo(x, y);
        context.lineTo(newX, newY);
        context.stroke();
        x = newX;
        y = newY;
      }
    }

    function drawLineTouch(event: TouchEvent) {
      if (isMouseDown) {
        const rect = canvasRef.current!.getBoundingClientRect()
        const touch = event.targetTouches[0];
        const newX = (touch.clientX - rect.x) / scale;
        const newY = (touch.clientY - rect.y) / scale;
        const context = canvasRef.current!.getContext("2d")!
        context.beginPath();
        context.moveTo(x, y);
        context.lineTo(newX, newY);
        context.stroke();
        x = newX;
        y = newY;
      }
    }

    function stopDrawing() {
      isMouseDown = false;
    }

    return <canvas ref={canvasRef}
                   onMouseDown={startDrawing}
                   onMouseMove={drawLine}
                   onMouseUp={stopDrawing}
                   onMouseOut={stopDrawing}
                   onTouchStart={startDrawingTouch}
                   onTouchMove={drawLineTouch}
                   onTouchEnd={stopDrawing}
                   onTouchCancel={stopDrawing}
                   width="820" height="1200"
                   style={{width: 820 * scale, height: 1200 * scale}}
    ></canvas>
  })

export default Canvas