const DATASETS = [
	'pages',
	'time',
	'referrers',
	'sources',
	'mediums',
	'campaigns',
	'browsers',
	'os',
	'devices',
	'countries',
	'languages',
] as const;

type Dataset = (typeof DATASETS)[number];

interface DataRow {
	// Common
	visitors?: number;
	visitors_percentage?: number;
	// Mixed
	path?: string;
	bounces?: number;
	bounce_rate?: number;
	duration?: number;
	// Pages
	pageviews?: number;
	pageviews_percentage?: number;
	// Duration
	duration_upper_quartile?: number;
	duration_lower_quartile?: number;
	duration_percentage?: number;
	// Referrers
	referrer?: string;
	// Sources
	source?: string;
	// Mediums
	medium?: string;
	// Campaigns
	campaign?: string;
	// Browsers
	browser?: string;
	// Operating Systems
	os?: string;
	// Devices
	device?: string;
	// Countries
	country?: string;
	// Languages
	language?: string;
}

interface StatsValue {
	label: string;
	count?: number;
	percentage?: number;
}

interface StatsTab {
	label: string;
	items: StatsValue[];
}

interface StatsGroups {
	label: string;
	data: StatsTab[];
}

export { DATASETS };
export type { DataRow, Dataset, StatsTab, StatsValue, StatsGroups };
