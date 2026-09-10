import{Tr as e,Zi as t,aa as n,bi as r,br as i,ia as a,or as o,rr as s,yi as c,yr as l}from"./router-omaRgp3e.js";import{n as u}from"./index-BXOIXutm.js";var d=c(`statistic`,[r(`label`,`
 font-weight: var(--n-label-font-weight);
 transition: .3s color var(--n-bezier);
 font-size: var(--n-label-font-size);
 color: var(--n-label-text-color);
 `),c(`statistic-value`,`
 margin-top: 4px;
 font-weight: var(--n-value-font-weight);
 `,[r(`prefix`,`
 margin: 0 4px 0 0;
 font-size: var(--n-value-font-size);
 transition: .3s color var(--n-bezier);
 color: var(--n-value-prefix-text-color);
 `,[c(`icon`,{verticalAlign:`-0.125em`})]),r(`content`,`
 font-size: var(--n-value-font-size);
 transition: .3s color var(--n-bezier);
 color: var(--n-value-text-color);
 `),r(`suffix`,`
 margin: 0 0 0 4px;
 font-size: var(--n-value-font-size);
 transition: .3s color var(--n-bezier);
 color: var(--n-value-suffix-text-color);
 `,[c(`icon`,{verticalAlign:`-0.125em`})])])]),f=a({name:`Statistic`,props:Object.assign(Object.assign({},s.props),{tabularNums:Boolean,label:String,value:[String,Number]}),slots:Object,setup(e){let{mergedClsPrefixRef:n,inlineThemeDisabled:r,mergedRtlRef:a}=i(e),c=s(`Statistic`,`-statistic`,d,u,e,n),f=o(`Statistic`,a,n),p=t(()=>{let{self:{labelFontWeight:e,valueFontSize:t,valueFontWeight:n,valuePrefixTextColor:r,labelTextColor:i,valueSuffixTextColor:a,valueTextColor:o,labelFontSize:s},common:{cubicBezierEaseInOut:l}}=c.value;return{"--n-bezier":l,"--n-label-font-size":s,"--n-label-font-weight":e,"--n-label-text-color":i,"--n-value-font-weight":n,"--n-value-font-size":t,"--n-value-prefix-text-color":r,"--n-value-suffix-text-color":a,"--n-value-text-color":o}}),m=r?l(`statistic`,void 0,p,e):void 0;return{rtlEnabled:f,mergedClsPrefix:n,cssVars:r?void 0:p,themeClass:m?.themeClass,onRender:m?.onRender}},render(){var t;let{mergedClsPrefix:r,$slots:{default:i,label:a,prefix:o,suffix:s}}=this;return(t=this.onRender)==null||t.call(this),n(`div`,{class:[`${r}-statistic`,this.themeClass,this.rtlEnabled&&`${r}-statistic--rtl`],style:this.cssVars},e(a,e=>n(`div`,{class:`${r}-statistic__label`},this.label||e)),n(`div`,{class:`${r}-statistic-value`,style:{fontVariantNumeric:this.tabularNums?`tabular-nums`:``}},e(o,e=>e&&n(`span`,{class:`${r}-statistic-value__prefix`},e)),this.value===void 0?e(i,e=>e&&n(`span`,{class:`${r}-statistic-value__content`},e)):n(`span`,{class:`${r}-statistic-value__content`},this.value),e(s,e=>e&&n(`span`,{class:`${r}-statistic-value__suffix`},e))))}});export{f as t};