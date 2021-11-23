import React from 'react';
import {queryByText, render, screen} from '@testing-library/react';
import LeftPanel from './LeftPanel';

it('renders correct header', () => {
  render(<LeftPanel />);
  expect(screen.getByTestId("left-panel-header")).toContainHTML("Open Pollution");
});
