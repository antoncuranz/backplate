import React, {useRef, useEffect, MouseEvent, TouchEvent, useImperativeHandle} from 'react'

interface Props {
  lineWidth: number,
  color: string,
}

const Canvas = React.forwardRef<HTMLCanvasElement, Props>(
  ({lineWidth, color}: Props, forwardedRef: React.ForwardedRef<HTMLCanvasElement>) => {

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

    function calculateScale() {
      return canvasRef.current!.offsetWidth / 820;
    }

    function startDrawing(event: MouseEvent) {
      isMouseDown = true;
      x = (event.nativeEvent.offsetX) / calculateScale();
      y = (event.nativeEvent.offsetY) / calculateScale();
    }

    function startDrawingTouch(event: TouchEvent) {
      isMouseDown = true;
      const rect = canvasRef.current!.getBoundingClientRect()
      const touch = event.targetTouches[0];
      x = (touch.clientX - rect.x) / calculateScale();
      y = (touch.clientY - rect.y) / calculateScale();
    }

    function drawLine(event: MouseEvent) {
      if (isMouseDown) {
        const newX = event.nativeEvent.offsetX / calculateScale();
        const newY = event.nativeEvent.offsetY / calculateScale();
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
        const newX = (touch.clientX - rect.x) / calculateScale();
        const newY = (touch.clientY - rect.y) / calculateScale();
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
    ></canvas>
  })

export default Canvas