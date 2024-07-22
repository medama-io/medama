export default {
	plugins: {
		'postcss-preset-mantine': {},
		'postcss-simple-vars': {
			variables: {
				'mantine-breakpoint-xs': '36em',
				'mantine-breakpoint-sm': '48em',
				'mantine-breakpoint-md': '62em',
				'mantine-breakpoint-lg': '75em',
				'mantine-breakpoint-xl': '88em',
			},
		},
		'postcss-lightningcss': {
			lightningcssOptions: {
				// Individually enable various drafts
				drafts: {
					// Enable custom media queries
					// https://drafts.csswg.org/mediaqueries-5/#custom-mq
					customMedia: true,
				},
			},
		},
	},
};
