import { Menu } from '@mantine/core';
import { ChevronDown } from 'lucide-react';
import { useMemo } from 'react';

import { ScrollArea } from '@/components/ScrollArea';

import classes from './WebsiteSelector.module.css';

interface WebsiteListComboboxProps {
	websites: string[];
	website: string;
	setWebsite: (website: string) => void;
}

export const WebsiteSelector = ({
	websites,
	website,
	setWebsite,
}: WebsiteListComboboxProps) => {
	const options = useMemo(
		() =>
			websites.map((value) => (
				<Menu.Item
					key={value}
					className={classes.item}
					data-active={value === website}
					onClick={() => setWebsite(value)}
				>
					{value}
				</Menu.Item>
			)),
		[websites, website, setWebsite],
	);

	return (
		<Menu
			position="bottom-start"
			offset={8}
			withinPortal
			classNames={{ dropdown: classes.dropdown }}
		>
			<Menu.Target>
				<button
					type="button"
					className={classes.trigger}
					aria-label="Select website hostname"
					disabled={websites.length === 0}
				>
					<span className={classes.label}>{website ?? 'No websites'}</span>
					<ChevronDown size={16} />
				</button>
			</Menu.Target>

			<Menu.Dropdown data-scroll={options.length > 5}>
				<ScrollArea vertical>{options}</ScrollArea>
			</Menu.Dropdown>
		</Menu>
	);
};
