export function authorNamesFilter(article) {
    let authors = [];

    for(let i = 0; i < article.latest_article_content.authors.length; i++) {
        authors.push(`${article.latest_article_content.authors[i].firstname} ${article.latest_article_content.authors[i].lastname}`);
    }

    return authors.join(", ");
}
