import React from 'react';

function LeftPanel() {
  return (
  <div className="LeftPanel">
    <h1 data-testid="left-panel-header">🏭 Open Pollution</h1>
    <hr className="MenuLine"/>
    <div className="SideNav">
      <a href="#">🗺️ Mapa pomiarów</a>
    </div>
  </div>

);
}

export default LeftPanel;
