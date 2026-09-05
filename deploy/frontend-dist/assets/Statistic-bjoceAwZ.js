import{Bi as e,Vn as t,Wn as n,Yi as r,ar as i,ir as a,oi as o,qi as s,si as c,ur as l}from"./router-DxH1-4bP.js";import{n as u}from"./index-BYJvH7k2.js";var d=o(`statistic`,[c(`label`,`
 font-weight: var(--n-label-font-weight);
 transition: .3s color var(--n-bezier);
 font-size: var(--n-label-font-size);
 color: var(--n-label-text-color);
 `),o(`statistic-value`,`
 margin-top: 4px;
 font-weight: var(--n-value-font-weight);
 `,[c(`prefix`,`
 margin: 0 4px 0 0;
 font-size: var(--n-value-font-size);
 transition: .3s color var(--n-bezier);
 color: var(--n-value-prefix-text-color);
 `,[o(`icon`,{verticalAlign:`-0.125em`})]),c(`content`,`
 font-size: var(--n-value-font-size);
 transition: .3s color var(--n-bezier);
 color: var(--n-value-text-color);
 `),c(`suffix`,`
 margin: 0 0 0 4px;
 font-size: var(--n-value-font-size);
 transition: .3s color var(--n-bezier);
 color: var(--n-value-suffix-text-color);
 `,[o(`icon`,{verticalAlign:`-0.125em`})])])]),f=s({name:`Statistic`,props:Object.assign(Object.assign({},t.props),{tabularNums:Boolean,label:String,value:[String,Number]}),slots:Object,setup(r){let{mergedClsPrefixRef:o,inlineThemeDisabled:s,mergedRtlRef:c}=i(r),l=t(`Statistic`,`-statistic`,d,u,r,o),f=n(`Statistic`,c,o),p=e(()=>{let{self:{labelFontWeight:e,valueFontSize:t,valueFontWeight:n,valuePrefixTextColor:r,labelTextColor:i,valueSuffixTextColor:a,valueTextColor:o,labelFontSize:s},common:{cubicBezierEaseInOut:c}}=l.value;return{"--n-bezier":c,"--n-label-font-size":s,"--n-label-font-weight":e,"--n-label-text-color":i,"--n-value-font-weight":n,"--n-value-font-size":t,"--n-value-prefix-text-color":r,"--n-value-suffix-text-color":a,"--n-value-text-color":o}}),m=s?a(`statistic`,void 0,p,r):void 0;return{rtlEnabled:f,mergedClsPrefix:o,cssVars:s?void 0:p,themeClass:m?.themeClass,onRender:m?.onRender}},render(){var e;let{mergedClsPrefix:t,$slots:{default:n,label:i,prefix:a,suffix:o}}=this;return(e=this.onRender)==null||e.call(this),r(`div`,{class:[`${t}-statistic`,this.themeClass,this.rtlEnabled&&`${t}-statistic--rtl`],style:this.cssVars},l(i,e=>r(`div`,{class:`${t}-statistic__label`},this.label||e)),r(`div`,{class:`${t}-statistic-value`,style:{fontVariantNumeric:this.tabularNums?`tabular-nums`:``}},l(a,e=>e&&r(`span`,{class:`${t}-statistic-value__prefix`},e)),this.value===void 0?l(n,e=>e&&r(`span`,{class:`${t}-statistic-value__content`},e)):r(`span`,{class:`${t}-statistic-value__content`},this.value),l(o,e=>e&&r(`span`,{class:`${t}-statistic-value__suffix`},e))))}});export{f as t};