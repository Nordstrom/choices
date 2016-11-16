import React from 'react';

const Segment = ({ namespaceSegments, experimentSegments = [], cellSize, spacing, rows }) => {
  const cols = Math.ceil(128 / rows);
  const segs = new Array(128).fill(false);
  namespaceSegments.forEach((s, i) => {
    if (s === 0) {
      return;
    }
    segs[i] = <rect
      key={i}
      fill="#00295b"
      width={cellSize}
      height={cellSize}
      x={Math.floor(i/rows) * (cellSize + spacing)}
      y={(i%rows) * (cellSize + spacing)}
    />;
  });
  experimentSegments.forEach((s, i) => {
    if (s === 0) {
      return;
    }
    segs[i] = <rect
      key={i}
      fill="#69bd28"
      width={cellSize}
      height={cellSize}
      x={Math.floor(i / rows) * (cellSize + spacing)}
      y={(i % rows) * (cellSize + spacing)}
    />
  });
  const out = segs.map((s, i) => {
    if (s) {
      return s;
    }
    return <rect
      key={i}
      fill="#b5b6ba"
      width={cellSize}
      height={cellSize}
      x={Math.floor(i/rows) * (cellSize + spacing)}
      y={(i%rows) * (cellSize + spacing)}
    />
  });
  return (
    <svg width={cellSize*cols + (cols-1)*spacing} height={cellSize*rows + (rows-1) * spacing}>
      {out}
    </svg>
  );
}

Segment.defaultProps = {
  namespaceSegments: [],
  experimentSegments: [],
  cellSize: 16,
  spacing: 2,
  rows: 4,
};

export default Segment;
