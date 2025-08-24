"use client";

import { cn } from "@/utils/utils";
import { GridPattern } from "@/components/magicui/grid-pattern";

export function GridPatternLinearGradient() {
  return (
    <div className="absolute inset-0 flex items-center justify-center z-0">
      <div className="relative flex size-full h-[500px] w-[1000px] items-center justify-center overflow-hidden rounded-lg bg-transparent p-20">
        <GridPattern
          width={64}
          height={64}
          x={-1}
          y={-1}
          className={cn(
            "[mask-image:radial-gradient(300px_circle_at_center,#eeeef2,transparent)]",
          )}
        />
      </div>
    </div>
  );
}