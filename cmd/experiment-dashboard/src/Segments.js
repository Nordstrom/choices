import React from 'react';
import "./Segments.css"

const off1 = <div className="off width-1"/>;
const off2 = <div className="off width-2"/>;
const off3 = <div className="off width-3"/>;
const off4 = <div className="off width-4"/>;

const on1 = <div className="on width-1"/>;
const on2 = <div className="on width-2"/>;
const on3 = <div className="on width-3"/>;
const on4 = <div className="on width-4"/>;

export const Segments = props => {
  const segments = props.segments.split('').map((char, i) => {
    switch (char) {
      case '0':
      return <div className="seg-half-byte" key={i}>{off4}</div>
      case '1':
      return <div className="seg-half-byte" key={i}>{off3}{on1}</div>
      case '2':
      return <div className="seg-half-byte" key={i}>{off2}{on1}{off1}</div>
      case '3':
      return <div className="seg-half-byte" key={i}>{off2}{on2}</div>
      case '4':
      return <div className="seg-half-byte" key={i}>{off1}{on1}{off2}</div>
      case '5':
      return <div className="seg-half-byte" key={i}>{off1}{on1}{off1}{on1}</div>
      case '6':
      return <div className="seg-half-byte" key={i}>{off1}{on2}{off1}</div>
      case '7':
      return <div className="seg-half-byte" key={i}>{off1}{on3}</div>
      case '8':
      return <div className="seg-half-byte" key={i}>{on1}{off3}</div>
      case '9':
      return <div className="seg-half-byte" key={i}>{on1}{off2}{on1}</div>
      case 'a':
      return <div className="seg-half-byte" key={i}>{on1}{off1}{on1}{off1}</div>
      case 'b':
      return <div className="seg-half-byte" key={i}>{on1}{off1}{on2}</div>
      case 'c':
      return <div className="seg-half-byte" key={i}>{on2}{off2}</div>
      case 'd':
      return <div className="seg-half-byte" key={i}>{on2}{off1}{on1}</div>
      case 'e':
      return <div className="seg-half-byte" key={i}>{on3}{off1}</div>
      case 'f':
      return <div className="seg-half-byte" key={i}>{on4}</div>
      default:
      return <div className="seg-half-byte" key={i}>{off4}</div>
    }
  });
  return (
    <div className="segments">
      {segments}
    </div>
  );
}