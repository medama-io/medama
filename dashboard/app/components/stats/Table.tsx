import { ChevronLeftIcon, ChevronRightIcon } from '@radix-ui/react-icons';
import * as Tabs from '@radix-ui/react-tabs';
import { useNavigate, useSearchParams } from '@remix-run/react';
import {
	DataTable,
	type DataTableColumn,
	type DataTableRowClickHandler,
	type DataTableSortStatus,
} from 'mantine-datatable';
import { useCallback, useMemo, useState } from 'react';
import isFQDN from 'validator/lib/isFQDN';

import { ButtonIcon, ButtonLink } from '@/components/Button';
import { IconExternal } from '@/components/icons/external';
import { Group } from '@/components/layout/Flex';
import { useDidUpdate } from '@/hooks/use-did-update';
import { useFilter } from '@/hooks/use-filter';
import { useHover } from '@/hooks/use-hover';
import { useMediaQuery } from '@/hooks/use-media-query';

import { formatCount, formatDuration, formatPercentage } from './formatter';
import type { DataRow, Dataset, Filter } from './types';
import { sortBy } from './utils';

import classes from './Table.module.css';

type DataRowClick = DataTableRowClickHandler<DataRow>;

interface TablePaginationProps {
	page: number;
	pageSize: number;
	totalRecords: number;
	onPageChange: (page: number) => void;
	onPageSizeChange: (pageSize: (typeof PAGE_SIZES)[number]) => void;
}

interface StatsTableProps {
	query: Dataset;
	data: DataRow[];
}

interface QueryTableProps extends StatsTableProps {
	isMobile?: boolean;
}

const LABEL_MAP: Record<Dataset, string> = {
	pages: 'Pages',
	time: 'Time Spent',
	referrers: 'Referrers',
	sources: 'Sources',
	mediums: 'Mediums',
	campaigns: 'Campaigns',
	browsers: 'Browsers',
	os: 'Operating Systems',
	devices: 'Devices',
	countries: 'Countries',
	languages: 'Languages',
	properties: 'Properties',
} as const;

const ACCESSOR_MAP: Record<Dataset, keyof DataRow> = {
	pages: 'path',
	time: 'duration',
	referrers: 'referrer',
	sources: 'source',
	mediums: 'medium',
	campaigns: 'campaign',
	browsers: 'browser',
	os: 'os',
	devices: 'device',
	countries: 'country',
	languages: 'language',
	properties: 'name',
} as const;

const FILTER_MAP: Record<Dataset, Filter> = {
	pages: 'path',
	time: 'path',
	referrers: 'referrer',
	sources: 'utm_source',
	mediums: 'utm_medium',
	campaigns: 'utm_campaign',
	browsers: 'browser',
	os: 'os',
	devices: 'device',
	countries: 'country',
	languages: 'language',
	properties: 'prop_name',
} as const;

const PAGE_SIZES = [10, 25, 50, 100] as const;

// Preset columns
type PresetDataKeys =
	| 'path'
	| 'visitors'
	| 'visitors_percentage'
	| 'pageviews'
	| 'pageviews_percentage'
	| 'bounce_percentage'
	| 'duration'
	| 'duration_lower_quartile'
	| 'duration_upper_quartile'
	| 'duration_percentage';

const PRESET_COLUMNS: Record<PresetDataKeys, DataTableColumn<DataRow>> = {
	path: { accessor: 'path', title: 'Path', width: '100%' },
	visitors: {
		accessor: 'visitors',
		title: 'Visitors',
		sortable: true,
		render: (record) => formatCount(record.visitors ?? 0),
	},
	visitors_percentage: {
		accessor: 'visitors_percentage',
		title: 'Visitors %',
		render: (record) => formatPercentage(record.visitors_percentage ?? 0),
	},
	pageviews: {
		accessor: 'pageviews',
		title: 'Views',
		sortable: true,
		render: (record) => formatCount(record.pageviews ?? 0),
	},
	pageviews_percentage: {
		accessor: 'pageviews_percentage',
		title: 'Views %',
		render: (record) => formatPercentage(record.pageviews_percentage ?? 0),
	},
	bounce_percentage: {
		accessor: 'bounce_percentage',
		title: 'Bounce %',
		render: (record) => formatPercentage(record.bounce_percentage ?? 0),
	},
	duration: {
		accessor: 'duration',
		title: 'Duration',
		sortable: true,
		render: (record) => formatDuration(record.duration ?? 0),
	},
	duration_lower_quartile: {
		accessor: 'duration_lower_quartile',
		title: 'Q1 (25%)',
		sortable: true,
		render: (record) => formatDuration(record.duration_lower_quartile ?? 0),
	},
	duration_upper_quartile: {
		accessor: 'duration_upper_quartile',
		title: 'Q3 (75%)',
		sortable: true,
		render: (record) => formatDuration(record.duration_upper_quartile ?? 0),
	},
	duration_percentage: {
		accessor: 'duration_percentage',
		title: 'Duration %',
		sortable: true,
		render: (record) => formatPercentage(record.duration_percentage ?? 0),
	},
};

const BackButton = () => {
	const [searchParams] = useSearchParams();

	return (
		<ButtonLink
			to={{
				pathname: '../',
				search: searchParams.toString(),
			}}
			className={classes.back}
		>
			<ChevronLeftIcon />
			<span>Go back</span>
		</ButtonLink>
	);
};

const QueryTable = ({ query, data, isMobile }: QueryTableProps) => {
	// Pagination
	const [pageSize, setPageSize] = useState<(typeof PAGE_SIZES)[number]>(10);
	const [page, setPage] = useState(1);

	// Sorting
	const [sortStatus, setSortStatus] = useState<DataTableSortStatus<DataRow>>({
		columnAccessor: query === 'properties' ? 'events' : 'visitors',
		direction: 'desc',
	});

	// Reset sort status when query changes.
	useDidUpdate(() => {
		setSortStatus({
			columnAccessor: query === 'properties' ? 'events' : 'visitors',
			direction: 'desc',
		});
	}, [query]);

	const { isFilterActiveEq } = useFilter();
	const isActiveFilter =
		isFilterActiveEq('language') || isFilterActiveEq('prop_name');
	const columns = useMemo(
		() => getColumnsForQuery(query, isActiveFilter),
		[query, isActiveFilter],
	);

	const records = useMemo(() => {
		const from = (page - 1) * pageSize;
		const to = from + pageSize;

		const sortedData = [...data].sort(
			sortBy(sortStatus.columnAccessor as keyof DataRow),
		);
		return sortStatus.direction === 'desc'
			? sortedData.reverse().slice(from, to)
			: sortedData.slice(from, to);
	}, [data, page, pageSize, sortStatus]);

	const handlePageChange = useCallback(
		(newPage: number) => {
			const maxPage = Math.ceil(data.length / pageSize);
			setPage(Math.max(1, Math.min(newPage, maxPage)));
		},
		[data.length, pageSize],
	);

	const handlePageSizeChange = useCallback(
		(newSize: (typeof PAGE_SIZES)[number]) => {
			setPageSize(newSize);
			setPage(1);
		},
		[],
	);

	// Reset page when query or data length changes (from filters).
	useDidUpdate(() => {
		setPage(1);
		setPageSize(10);
	}, [query, data.length]);

	// Filter
	const { addFilter } = useFilter();

	const handleFilter: DataRowClick = (row) => {
		const { record } = row;
		let filter = FILTER_MAP[query] ?? 'path';

		let value = record[ACCESSOR_MAP[query]] || 'Direct/None';
		// Time query uses path as the accessor and main sort.
		if (query === 'time' && record.path) {
			value = record.path;
			// If value exists, prop_name is set and we should access values now.
		} else if (query === 'properties' && record.value) {
			filter = 'prop_value';
			value = record.value;
		}

		addFilter(filter, 'eq', String(value));
	};

	return (
		<div className={classes['table-wrapper']}>
			<div className={classes['table-header']}>
				<span>{LABEL_MAP[query]}</span>
			</div>
			{isMobile && <BackButton />}
			<DataTable
				classNames={{ header: classes['data-header'] }}
				minHeight={365}
				noRecordsText="No records found..."
				highlightOnHover
				withRowBorders={false}
				idAccessor={(record) =>
					String(
						record[ACCESSOR_MAP[query]] ??
							record.name ??
							record.value ??
							'path',
					)
				}
				records={records}
				columns={columns}
				sortStatus={sortStatus}
				onSortStatusChange={setSortStatus}
				onRowClick={handleFilter}
			/>
			<TablePagination
				page={page}
				pageSize={pageSize}
				totalRecords={data.length}
				onPageChange={handlePageChange}
				onPageSizeChange={handlePageSizeChange}
			/>
		</div>
	);
};

const TablePagination = ({
	page,
	pageSize,
	totalRecords,
	onPageChange,
	onPageSizeChange,
}: TablePaginationProps) => {
	const totalPages = Math.ceil(totalRecords / pageSize);

	return (
		<div className={classes.pagination}>
			<Group data-visible-from="sm" style={{ gap: 16 }}>
				<span className={classes.viewspan}>View</span>
				{PAGE_SIZES.map((size) => (
					<ButtonIcon
						key={size}
						label={`Show ${size} records`}
						className={classes['page-size']}
						onClick={() => onPageSizeChange(size)}
						disabled={size === pageSize || totalRecords <= size}
						data-active={size === pageSize}
					>
						{size}
					</ButtonIcon>
				))}
			</Group>
			<Group style={{ gap: 16 }}>
				<ButtonIcon
					label="Previous page"
					className={classes['page-arrow']}
					onClick={() => onPageChange(page - 1)}
					disabled={page <= 1}
				>
					<ChevronLeftIcon />
				</ButtonIcon>
				<span>
					Page {page} of {totalPages}
				</span>
				<ButtonIcon
					label="Next page"
					className={classes['page-arrow']}
					onClick={() => onPageChange(page + 1)}
					disabled={page >= totalPages}
				>
					<ChevronRightIcon />
				</ButtonIcon>
			</Group>
		</div>
	);
};

const getColumnsForQuery = (
	query: Dataset,
	filterActive?: boolean,
): DataTableColumn<DataRow>[] => {
	switch (query) {
		case 'pages':
			return [
				PRESET_COLUMNS.path,
				PRESET_COLUMNS.visitors,
				PRESET_COLUMNS.visitors_percentage,
				PRESET_COLUMNS.pageviews,
				PRESET_COLUMNS.pageviews_percentage,
				PRESET_COLUMNS.bounce_percentage,
				PRESET_COLUMNS.duration,
			];
		case 'time':
			return [
				PRESET_COLUMNS.path,
				PRESET_COLUMNS.visitors,
				PRESET_COLUMNS.duration,
				PRESET_COLUMNS.duration_lower_quartile,
				PRESET_COLUMNS.duration_upper_quartile,
				PRESET_COLUMNS.duration_percentage,
			];
		case 'referrers':
			return [
				{
					// Referrer, Source, Medium, Campaign
					accessor: ACCESSOR_MAP[query],
					title: LABEL_MAP[query].slice(0, -1),
					width: '100%',
					render: (record) => {
						const label = record[ACCESSOR_MAP[query]];
						const { ref, hovered } = useHover<HTMLDivElement>();

						if (label) {
							return (
								<div ref={ref}>
									{label}
									<a
										className={classes.external}
										aria-label={`Visit ${label}`}
										href={`https://${label}`}
										target="_blank"
										rel="noreferrer noopener"
										data-hidden={!isFQDN(String(label))}
										data-hover={hovered ? 'true' : undefined}
										onClick={(event) => event.stopPropagation()}
									>
										<IconExternal />
									</a>
								</div>
							);
						}

						return 'Direct/None';
					},
				},
				PRESET_COLUMNS.visitors,
				PRESET_COLUMNS.visitors_percentage,
				PRESET_COLUMNS.bounce_percentage,
				PRESET_COLUMNS.duration,
			];
		case 'sources':
		case 'mediums':
		case 'campaigns':
			return [
				{
					// Referrer, Source, Medium, Campaign
					accessor: ACCESSOR_MAP[query],
					title: LABEL_MAP[query].slice(0, -1),
					width: '100%',
					render: (record) => record[ACCESSOR_MAP[query]] || 'Direct/None',
				},
				PRESET_COLUMNS.visitors,
				PRESET_COLUMNS.visitors_percentage,
				PRESET_COLUMNS.bounce_percentage,
				PRESET_COLUMNS.duration,
			];
		case 'browsers':
		case 'devices':
			return [
				{
					// Browser, Device
					accessor: ACCESSOR_MAP[query],
					title: LABEL_MAP[query].slice(0, -1),
					width: '100%',
					render: (record) => record[ACCESSOR_MAP[query]] || 'Unknown',
				},
				PRESET_COLUMNS.visitors,
				PRESET_COLUMNS.visitors_percentage,
				PRESET_COLUMNS.bounce_percentage,
				PRESET_COLUMNS.duration,
			];
		case 'languages':
			return [
				{
					// Browser, Device, Language
					accessor: ACCESSOR_MAP[query],
					title: filterActive ? 'Locale' : 'Language',
					width: '100%',
					render: (record) => record[ACCESSOR_MAP[query]] || 'Unknown',
				},
				PRESET_COLUMNS.visitors,
				PRESET_COLUMNS.visitors_percentage,
				PRESET_COLUMNS.bounce_percentage,
				PRESET_COLUMNS.duration,
			];
		case 'os':
			return [
				{
					accessor: ACCESSOR_MAP[query],
					title: 'OS',
					width: '100%',
					render: (record) => record.os || 'Unknown',
				},
				PRESET_COLUMNS.visitors,
				PRESET_COLUMNS.visitors_percentage,
				PRESET_COLUMNS.bounce_percentage,
				PRESET_COLUMNS.duration,
			];
		case 'countries':
			return [
				{
					accessor: ACCESSOR_MAP[query],
					title: 'Country',
					width: '100%',
					render: (record) => record.country || 'Unknown',
				},
				PRESET_COLUMNS.visitors,
				PRESET_COLUMNS.visitors_percentage,
				PRESET_COLUMNS.bounce_percentage,
				PRESET_COLUMNS.duration,
			];
		case 'properties':
			return [
				{
					accessor: filterActive ? 'value' : 'name',
					title: 'Name',
					width: '100%',
					render: (record) => record.name || record.value || 'Unknown',
				},
				{
					accessor: 'events',
					title: 'Events',
					sortable: true,
					render: (record) => formatCount(record.events ?? 0),
				},
				{
					accessor: 'events_percentage',
					title: 'Events %',
					sortable: true,
					render: (record) => formatPercentage(record.events_percentage ?? 0),
				},
			];
		default:
			throw new Error(`Invalid query: ${query}`);
	}
};

export const Table = ({ query, data }: StatsTableProps) => {
	const navigate = useNavigate();
	const [searchParams] = useSearchParams();

	const handleTabChange = useCallback(
		(value: string | null) => {
			navigate(
				{
					pathname: `../${value}`,
					search: searchParams.toString(),
				},
				{ preventScrollReset: true },
			);
		},
		[navigate, searchParams],
	);

	const isMobile = useMediaQuery('(max-width: 75em)');

	return (
		<Tabs.Root
			value={query}
			className={classes.root}
			orientation="vertical"
			onValueChange={handleTabChange}
		>
			{!isMobile && (
				<div className={classes.list}>
					<BackButton />
					<Tabs.List className={classes['list-triggers']}>
						{Object.entries(LABEL_MAP).map(([key, label]) => (
							<Tabs.Trigger key={key} value={key} className={classes.trigger}>
								{label}
							</Tabs.Trigger>
						))}
					</Tabs.List>
				</div>
			)}
			<Tabs.Content value={query} className={classes.panel}>
				<QueryTable query={query} data={data} isMobile={isMobile} />
			</Tabs.Content>
		</Tabs.Root>
	);
};
