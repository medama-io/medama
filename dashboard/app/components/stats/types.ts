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
	'properties',
] as const;

type Dataset = (typeof DATASETS)[number];

interface DataRow {
	// Common
	visitors?: number;
	visitors_percentage?: number;
	// Mixed
	path?: string;
	bounce_percentage?: number;
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
	// Properties
	property?: string;
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

// Charts

type ChartStat = 'visitors' | 'pageviews' | 'duration' | 'bounces';
type ChartType = 'area' | 'bar';

interface StatHeaderData {
	label: string;
	chart: ChartStat;
	current: number;
	previous?: number;
}

// Filters

const FILTERS = [
	'path',
	'referrer',
	'utm_source',
	'utm_medium',
	'utm_campaign',
	'browser',
	'os',
	'device',
	'country',
	'language',
] as const;

type Filter = (typeof FILTERS)[number];

const FILTER_OPERATOR = [
	'eq',
	'neq',
	'contains',
	'not_contains',
	'starts_with',
	'not_starts_with',
	'ends_with',
	'not_ends_with',
] as const;

type FilterOperator = (typeof FILTER_OPERATOR)[number];

type FilterKey = `${Filter}[${FilterOperator}]`;

export { DATASETS };
export type {
	ChartStat,
	ChartType,
	DataRow,
	Dataset,
	StatHeaderData,
	StatsGroups,
	StatsTab,
	StatsValue,
	Filter,
	FilterOperator,
	FilterKey,
};
