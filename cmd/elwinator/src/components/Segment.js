import React from 'react';

const cellSize = 18;

const Segment = ({ namespaceSegments, experimentSegments = [] }) => {
  const segs = new Array(128).fill(false);
  namespaceSegments.forEach((s, i) => s === 0 ? false : 
    segs[i] = <rect key={i} fill="#00295b" width={cellSize-2} height={cellSize-2} x={(Math.floor(i/4) * cellSize + 1)} y={(i%4) * cellSize + 1} />
  );
  experimentSegments.forEach((s, i) => s === 0 ? false :
    segs[i] = <rect key={i} fill="#69bd28" width={cellSize-2} height={cellSize-2} x={(Math.floor(i / 4) * cellSize + 1) } y={(i % 4) * cellSize + 1} />
  );
  const out = segs.map((s, i) => {
    if (s) {
      return s;
    }
    return <rect key={i} fill="#b5b6ba" width={cellSize-2} height={cellSize-2} x={(Math.floor(i/4) * cellSize + 1)} y={(i%4) * cellSize + 1} />
  });
  return (
    <svg width={cellSize*32} height={cellSize*4}>
      {out}
    </svg>
  );
}

export default Segment;
