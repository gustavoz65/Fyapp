-- Fix category icons to use emojis instead of icon names
UPDATE categories SET icon = '🍴' WHERE icon = 'utensils';
UPDATE categories SET icon = '🚗' WHERE icon = 'car';
UPDATE categories SET icon = '🏠' WHERE icon = 'home';
UPDATE categories SET icon = '❤️' WHERE icon = 'heart';
UPDATE categories SET icon = '📚' WHERE icon = 'book';
UPDATE categories SET icon = '🎮' WHERE icon = 'gamepad';
UPDATE categories SET icon = '🛍️' WHERE icon = 'shopping-bag';
UPDATE categories SET icon = '⚙️' WHERE icon = 'settings';
UPDATE categories SET icon = '➕' WHERE icon = 'more-horizontal';
UPDATE categories SET icon = '💼' WHERE icon = 'briefcase';
UPDATE categories SET icon = '💻' WHERE icon = 'laptop';
UPDATE categories SET icon = '📈' WHERE icon = 'trending-up';
UPDATE categories SET icon = '🎁' WHERE icon = 'gift';
UPDATE categories SET icon = '🔄' WHERE icon = 'refresh-cw';
UPDATE categories SET icon = '💰' WHERE icon = 'wallet';
UPDATE categories SET icon = '🏦' WHERE icon = 'bank';
