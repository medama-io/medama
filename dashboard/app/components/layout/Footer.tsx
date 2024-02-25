import { Box, Text } from '@mantine/core';

import classes from './Footer.module.css';

export const Footer = () => {
	return (
		<Box component="footer" className={classes.footer}>
			<div className={classes.inner}>
				<Text></Text>
			</div>
		</Box>
	);
};
