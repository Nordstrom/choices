import React from 'react';

const NavSection = ({ children }) => {
  const links = children.map((node, i) => <li key={i} className="nav-item">{node}</li>);
  return (
    <div className="row">
      <div className="col-sm-3">
        <nav>
          <ul className="nav nav-pills nav-stacked">
            {links}
          </ul>
        </nav>
      </div>
    </div>
  );
}

export default NavSection;
