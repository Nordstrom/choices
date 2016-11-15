import React from 'react';
import { Link } from 'react-router';

import { rootURL } from '../urls';
import Back from './Back';

const links = (children = []) => {
  let ls;
  const home = <Link to={ rootURL() }>Home</Link>;
  if (children.constructor === Array) {
    ls = [home, ...children, <Back />];
  } else if (!children) {
    ls =[home, <Back />];
  } else {
    ls = [home, children, <Back />];
  }
  return ls.map((l, i) => <li key={i} className="nav-item">{l}</li>);
}

const NavSection = ({ children }) => {
  return (
      <div className="col-sm-3">
        <nav>
          <ul className="nav nav-pills nav-stacked">
            {links(children)}
          </ul>
        </nav>
      </div>
  );
}

export default NavSection;
