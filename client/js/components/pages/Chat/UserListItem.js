import React, { memo } from 'react';
import stringToRGB from 'utils/color';

const UserListItem = ({ user, coloredNick, style, onClick }) => {
  if (coloredNick) {
    style = {
      ...style,
      color: stringToRGB(user.nick)
    };
  }

  return (
    <p style={style} onClick={() => onClick(user.nick)}>
      {user.renderName}
    </p>
  );
};

export default memo(UserListItem);
